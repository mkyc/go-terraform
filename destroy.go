package terra

// Destroy runs terraform destroy with the given options and return stdout/stderr.
func Destroy(options *Options) (string, error) {
	return RunTerraformCommandE(options, FormatArgs(options, "destroy", "-auto-approve", "-input=false")...)
}
