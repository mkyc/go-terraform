package terra

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Command is a simpler struct for defining commands than Go's built-in Cmd.
type Command struct {
	Command    string            // The command to run
	Args       []string          // The args to pass to the command
	WorkingDir string            // The working directory
	Env        map[string]string // Additional environment variables to set
	// Use the specified logger for the command's output. Use logger.Discard to not print the output while executing the command.
	Logger Logger
}

// RunTerraformCommandE runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommandE(additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)
	return DoWithRetryableErrorsE(description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, options.Logger, func() (string, error) {
		return RunCommandAndGetOutputE(cmd)
	})
}

// RunTerraformCommandAndGetStdoutE runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
func RunTerraformCommandAndGetStdoutE(additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)
	return DoWithRetryableErrorsE(description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, options.Logger, func() (string, error) {
		return RunCommandAndGetStdOutE(cmd)
	})
}

func generateCommand(options *Options, args ...string) Command {
	cmd := Command{
		Command:    options.TerraformBinary,
		Args:       args,
		WorkingDir: options.TerraformDir,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
	return cmd
}

// RunCommandAndGetOutputE runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of
// that command will also be logged with Command.Log to make debugging easier. Any returned error will be of type
// ErrWithCmdOutput, containing the output streams and the underlying error.
func RunCommandAndGetOutputE(command Command) (string, error) {
	output, err := runCommand(command)
	if err != nil {
		if output != nil {
			return output.Combined(), &ErrWithCmdOutput{err, output}
		}
		return "", &ErrWithCmdOutput{err, output}
	}

	return output.Combined(), nil
}

// RunCommandAndGetStdOutE runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
// Any returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
func RunCommandAndGetStdOutE(command Command) (string, error) {
	output, err := runCommand(command)
	if err != nil {
		if output != nil {
			return output.Stdout(), &ErrWithCmdOutput{err, output}
		}
		return "", &ErrWithCmdOutput{err, output}
	}

	return output.Stdout(), nil
}

// runCommand runs a shell command and stores each line from stdout and stderr in Output. Depending on the logger, the
// stdout and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging
// easier.
func runCommand(command Command) (*output, error) {
	command.Logger.Info("Running command %s with args %s\n", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Env = formatEnvVars(command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	output, err := readStdoutAndStderr(command, stdout, stderr)
	if err != nil {
		return output, err
	}

	return output, cmd.Wait()
}

// This function captures stdout and stderr into the given variables while still printing it to the stdout and stderr
// of this Go program
func readStdoutAndStderr(command Command, stdout, stderr io.ReadCloser) (*output, error) {
	out := newOutput()
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	var stdoutErr, stderrErr error
	go func() {
		defer wg.Done()
		stdoutErr = readData(command, false, stdoutReader, out.stdout)
	}()
	go func() {
		defer wg.Done()
		stderrErr = readData(command, true, stderrReader, out.stderr)
	}()
	wg.Wait()

	if stdoutErr != nil {
		return out, stdoutErr
	}
	if stderrErr != nil {
		return out, stderrErr
	}

	return out, nil
}

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

func readData(command Command, isStderr bool, reader *bufio.Reader, writer io.StringWriter) error {
	var line string
	var readErr error
	for {
		line, readErr = reader.ReadString('\n')

		// remove newline, our output is in a slice,
		// one element per line.
		line = strings.TrimSuffix(line, "\n")

		// only return early if the line does not have
		// any contents. We could have a line that does
		// not not have a newline before io.EOF, we still
		// need to add it to the output.
		if len(line) == 0 && readErr == io.EOF {
			break
		}

		if isStderr {
			command.Logger.Warn(line)
		} else {
			command.Logger.Info(line)
		}

		if _, err := writer.WriteString(line); err != nil {
			return err
		}

		if readErr != nil {
			break
		}
	}
	if readErr != io.EOF {
		return readErr
	}
	return nil
}

// DoWithRetryableErrorsE runs the specified action. If it returns a value, return that value. If it returns an error,
// check if error message or the string output from the action (which is often stdout/stderr from running some command)
// matches any of the regular expressions in the specified retryableErrors map. If there is a match, sleep for
// sleepBetweenRetries, and retry the specified action, up to a maximum of maxRetries retries. If there is no match,
// return that error immediately, wrapped in a FatalError. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryableErrorsE(actionDescription string, retryableErrors map[string]string, maxRetries int, sleepBetweenRetries time.Duration, logger Logger, action func() (string, error)) (string, error) {
	retryableErrorsRegexp := map[*regexp.Regexp]string{}
	for errorStr, errorMessage := range retryableErrors {
		errorRegex, err := regexp.Compile(errorStr)
		if err != nil {
			return "", FatalError{Underlying: err}
		}
		retryableErrorsRegexp[errorRegex] = errorMessage
	}

	return DoWithRetryE(actionDescription, maxRetries, sleepBetweenRetries, logger, func() (string, error) {
		output, err := action()
		if err == nil {
			return output, nil
		}

		for errorRegexp, errorMessage := range retryableErrorsRegexp {
			if errorRegexp.MatchString(output) || errorRegexp.MatchString(err.Error()) {
				logger.Warn("'%s' failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", actionDescription, err.Error(), errorMessage)
				return output, err
			}
		}

		return output, FatalError{Underlying: err}
	})
}

// DoWithRetryE runs the specified action. If it returns a string, return that string. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryE(actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, logger Logger, action func() (string, error)) (string, error) {
	out, err := DoWithRetryInterfaceE(actionDescription, maxRetries, sleepBetweenRetries, logger, func() (interface{}, error) { return action() })
	return out.(string), err
}

// DoWithRetryInterfaceE runs the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryInterfaceE(actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, logger Logger, action func() (interface{}, error)) (interface{}, error) {
	var output interface{}
	var err error

	for i := 0; i <= maxRetries; i++ {
		logger.Info(actionDescription)

		output, err = action()
		if err == nil {
			return output, nil
		}

		if _, isFatalErr := err.(FatalError); isFatalErr {
			logger.Error("Returning due to fatal error: %v\n", err)
			return output, err
		}

		logger.Warn("%s returned an error: %s. Sleeping for %s and will try again.\n", actionDescription, err.Error(), sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}
