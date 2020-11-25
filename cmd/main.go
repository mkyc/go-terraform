package main

import (
	"fmt"
	terra "github.com/mkyc/go-terraform"
)

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
