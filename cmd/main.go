package main

import (
	"fmt"
	terra "github.com/mkyc/go-terraform"
)

type SimpleLogger struct{}

func (s SimpleLogger) Trace(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Debug(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Info(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Warn(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Error(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Fatal(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Panic(format string, v ...interface{}) {
	justPrint(format, v...)
}

func justPrint(s string, v ...interface{}) {
	if len(v) > 0 {
		fmt.Printf(s, v...)
	} else {
		fmt.Println(s)
	}
}

func main() {
	opts, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir:  "./tests",
		StateFilePath: "./subdir/other-state.tfstate",
		PlanFilePath:  "./subdir/other-plan.tfplan",
	})
	if err != nil {
		panic(err)
	}

	println("================")
	println("================")
	println("=== init =======")
	println("================")
	println("================")

	s, err := terra.Init(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("================")
	println("=== plan =======")
	println("================")
	println("================")

	s, err = terra.Plan(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("================")
	println("=== show =======")
	println("================")
	println("================")

	s, err = terra.Show(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("================")
	println("=== apply ======")
	println("================")
	println("================")

	s, err = terra.Apply(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("================")
	println("=== output =====")
	println("================")
	println("================")

	m, err := terra.OutputAll(opts)
	if err != nil {
		panic(err)
	}
	println(len(m))
	for k, v := range m {
		fmt.Printf("%s : %v\n", k, v)
	}

	println("================")
	println("================")
	println("=== destroy ====")
	println("================")
	println("================")

	s, err = terra.Destroy(opts)
	if err != nil {
		panic(err)
	}

	println("================")
	println("=== plan =======")
	println("================")
	println("=== apply ======")
	println("================")

	_, err = terra.Plan(opts)
	if err != nil {
		panic(err)
	}
	_, err = terra.Apply(opts)
	if err != nil {
		panic(err)
	}

	opts, err = terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir:  "./tests",
		StateFilePath: "./subdir/other-state.tfstate",
		PlanFilePath:  "./subdir/other-destroy-plan.tfplan",
		Logger:        SimpleLogger{},
	})
	if err != nil {
		panic(err)
	}

	println("================")
	println("================")
	println("= plan destroy =")
	println("================")
	println("================")

	s, err = terra.PlanDestroy(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("=== show =======")
	println("================")
	println("= destruction ==")
	println("================")

	s, err = terra.Show(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("================")
	println("= destroy via ==")
	println("=== plan of ====")
	println("= destruction===")
	println("================")

	s, err = terra.Apply(opts)
	if err != nil {
		panic(err)
	}
	println(s)
}
