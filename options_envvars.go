package god

type environmentVariablesOption interface {
	WithEnvironmentVariables(map[string]string) Options
	EnvironmentVariables() map[string]string
	HasEnvironmentVariables() bool
}

var _ environmentVariablesOption = &options{}

func (opts *options) WithEnvironmentVariables(envs map[string]string) Options {
	newOpts := opts.copy()

	environmentVariables := make(map[string]string, len(envs))
	for k, v := range envs {
		environmentVariables[k] = v
	}

	newOpts.environmentVariables = &environmentVariables
	return newOpts
}

func (opts *options) EnvironmentVariables() map[string]string {
	return *opts.environmentVariables
}

func (opts *options) HasEnvironmentVariables() bool {
	return opts.environmentVariables != nil
}
