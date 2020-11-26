package terra

import "log"

type Logger interface {

	// Trace arguments are handled in the manner of fmt.Printf.
	Trace(format string, v ...interface{})

	// Debug arguments are handled in the manner of fmt.Printf.
	Debug(format string, v ...interface{})

	// Info arguments are handled in the manner of fmt.Printf.
	Info(format string, v ...interface{})

	// Warn arguments are handled in the manner of fmt.Printf.
	Warn(format string, v ...interface{})

	// Error arguments are handled in the manner of fmt.Printf.
	Error(format string, v ...interface{})

	// Fatal arguments are handled in the manner of fmt.Printf.
	Fatal(format string, v ...interface{})

	// Panic arguments are handled in the manner of fmt.Printf.
	Panic(format string, v ...interface{})
}

type DefaultLogger struct {
	l Logger
}

// Trace calls log.Printf or log.Println (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Trace(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf("[Trace] "+format, v...)
	} else {
		log.Println("[Trace] " + format)
	}
}

// Debug calls log.Printf or log.Println (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Debug(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf("[Debug] "+format, v...)
	} else {
		log.Println("[Debug] " + format)
	}
}

// Info calls log.Printf or log.Println (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Info(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf("[Info] "+format, v...)
	} else {
		log.Println("[Info] " + format)
	}
}

// Warn calls log.Printf or log.Println (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Warn(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf("[Warn] "+format, v...)
	} else {
		log.Println("[Warn] " + format)
	}
}

// Error calls log.Printf or log.Println (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Error(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf("[Error] "+format, v...)
	} else {
		log.Println("[Error] " + format)
	}
}

// Fatal calls log.Fatalf or log.Fatalln (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Fatal(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Fatalf("[Trace] "+format, v...)
	} else {
		log.Fatalln("[Trace] " + format)
	}
}

// Panic calls log.Panicf or log.Panicln (if there is no
// additional arguments) to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (DefaultLogger) Panic(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Panicf("[Trace] "+format, v...)
	} else {
		log.Panicln("[Trace] " + format)
	}
}
