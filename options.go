package terra

import (
	"fmt"
	"github.com/jinzhu/copier"
	"time"
)

var (
	DefaultRetryableTerraformErrors = map[string]string{
		// Helm related terraform calls may fail when too many tests run in parallel. While the exact cause is unknown,
		// this is presumably due to all the network contention involved. Usually a retry resolves the issue.
		".*read: connection reset by peer.*": "Failed to reach helm charts repository.",
		".*transport is closing.*":           "Failed to reach Kubernetes API.",

		// `terraform init` frequently fails in CI due to network issues accessing plugins. The reason is unknown, but
		// eventually these succeed after a few retries.
		".*unable to verify signature.*":             "Failed to retrieve plugin due to transient network error.",
		".*unable to verify checksum.*":              "Failed to retrieve plugin due to transient network error.",
		".*no provider exists with the given name.*": "Failed to retrieve plugin due to transient network error.",
		".*registry service is unreachable.*":        "Failed to retrieve plugin due to transient network error.",
		".*Error installing provider.*":              "Failed to retrieve plugin due to transient network error.",

		// Provider bugs where the data after apply is not propagated. This is usually an eventual consistency issue, so
		// retrying should self resolve it.
		// See https://github.com/terraform-providers/terraform-provider-aws/issues/12449 for an example.
		".*Provider produced inconsistent result after apply.*": "Provider eventual consistency error.",
	}

	commandsWithParallelism = []string{
		"plan",
		"apply",
		"destroy",
		"plan-all",
		"apply-all",
		"destroy-all",
	}
)

// Options for running Terraform commands
type Options struct {
	//TODO remove because it is for terragrunt
	TerraformBinary string // Name of the binary that will be used
	TerraformDir    string // The path to the folder where the Terraform code is defined.

	// The vars to pass to Terraform commands using the -var option. Note that terraform does not support passing `null`
	// as a variable value through the command line. That is, if you use `map[string]interface{}{"foo": nil}` as `Vars`,
	// this will translate to the string literal `"null"` being assigned to the variable `foo`. However, nulls in
	// lists and maps/objects are supported. E.g., the following var will be set as expected (`{ bar = null }`:
	// map[string]interface{}{
	//     "foo": map[string]interface{}{"bar": nil},
	// }
	Vars map[string]interface{}

	VarFiles                 []string               // The var file paths to pass to Terraform commands using -var-file option.
	Targets                  []string               // The target resources to pass to the terraform command with -target
	Lock                     bool                   // The lock option to pass to the terraform command with -lock
	LockTimeout              string                 // The lock timeout option to pass to the terraform command with -lock-timeout
	EnvVars                  map[string]string      // Environment variables to set when running Terraform
	BackendConfig            map[string]interface{} // The vars to pass to the terraform init command for extra configuration for the backend
	RetryableTerraformErrors map[string]string      // If Terraform apply fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	MaxRetries               int                    // Maximum number of times to retry errors matching RetryableTerraformErrors
	TimeBetweenRetries       time.Duration          // The amount of time to wait between retries
	Upgrade                  bool                   // Whether the -upgrade flag of the terraform init command should be set to true or not
	NoColor                  bool                   // Whether the -no-color flag will be set for any Terraform command or not
	NoStderr                 bool                   // Disable stderr redirection
	OutputMaxLineSize        int                    // The max size of one line in stdout and stderr (in bytes)
	Parallelism              int                    // Set the parallelism setting for Terraform
	PlanFilePath             string                 // The path to output a plan file to (for the plan command) or read one from (for the apply command)
	StateFilePath            string                 // The path to state file
	Logger                   Logger
	//TODO consider adding ssh agent back
}

// Clone makes a deep copy of most fields on the Options object and returns it.
//
// NOTE: options.SshAgent and options.Logger CANNOT be deep copied (e.g., the SshAgent struct contains channels and
// listeners that can't be meaningfully copied), so the original values are retained.
func (options *Options) Clone() (*Options, error) {
	newOptions := &Options{}
	if err := copier.Copy(newOptions, options); err != nil {
		return nil, err
	}
	return newOptions, nil
}

// GetCommonOptions extracts commons terraform options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.NoColor && !listContains(args, "-no-color") {
		args = append(args, "-no-color")
	}

	if options.TerraformBinary == "" {
		options.TerraformBinary = "terraform"
	}

	if options.Parallelism > 0 && len(args) > 0 && listContains(commandsWithParallelism, args[0]) {
		args = append(args, fmt.Sprintf("--parallelism=%d", options.Parallelism))
	}

	if options.Logger == nil {
		options.Logger = &DefaultLogger{}
	}

	return options, args
}

// WithDefaultRetryableErrors makes a copy of the Options object and returns an updated object with sensible defaults
// for retryable errors. The included retryable errors are typical errors that most terraform modules encounter during
// testing, and are known to self resolve upon retrying.
// This will fail the test if there are any errors in the cloning process.
func WithDefaultRetryableErrors(originalOptions *Options) (*Options, error) {
	//TODO why cloning ... ?
	newOptions, err := originalOptions.Clone()
	if err != nil {
		return nil, err
	}

	if newOptions.RetryableTerraformErrors == nil {
		newOptions.RetryableTerraformErrors = map[string]string{}
	}
	for k, v := range DefaultRetryableTerraformErrors {
		newOptions.RetryableTerraformErrors[k] = v
	}

	// These defaults for retry configuration are arbitrary, but have worked well in practice across Gruntwork
	// modules.
	newOptions.MaxRetries = 3
	newOptions.TimeBetweenRetries = 5 * time.Second

	if newOptions.Logger == nil {
		newOptions.Logger = &DefaultLogger{}
	}

	return newOptions, nil
}
