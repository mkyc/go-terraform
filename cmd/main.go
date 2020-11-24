package main

import terra "github.com/mkyc/go-terraform"

func main() {
	opts, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: "./tests",
	})
	if err != nil {
		panic(err)
	}

	println("===============")
	println("===============")
	println("=== init ======")
	println("===============")
	println("===============")

	s, err := terra.Init(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("===============")
	println("===============")
	println("=== plan ======")
	println("===============")
	println("===============")

	s, err = terra.Plan(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("===============")
	println("===============")
	println("=== apply =====")
	println("===============")
	println("===============")

	s, err = terra.Apply(opts)
	if err != nil {
		panic(err)
	}
	println(s)

	println("===============")
	println("===============")
	println("=== destroy ===")
	println("===============")
	println("===============")

	s, err = terra.Destroy(opts)
	if err != nil {
		panic(err)
	}
}
