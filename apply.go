package terra

// Apply runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func Apply(options *Options) (string, error) {
	return RunTerraformCommandE(options, FormatArgs(options, "apply", "-input=false", "-auto-approve")...)
}
