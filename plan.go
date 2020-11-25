package terra

func Plan(options *Options) (string, error) {
	return RunTerraformCommandE(options, FormatArgs(options, "plan", "-input=false", "-lock=false")...)
}

func PlanDestroy(options *Options) (string, error) {
	return RunTerraformCommandE(options, FormatArgs(options, "plan", "-destroy", "-input=false", "-lock=false")...)
}
