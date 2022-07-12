package god

type argumentsOption interface {
	WithArguments(...string) Options
	Arguments() []string
	HasArguments() bool
}

var _ argumentsOption = &options{}

func (opts *options) WithArguments(arguments ...string) Options {
	newOpts := opts.copy()
	newOpts.arguments = &arguments
	return newOpts
}

func (opts *options) Arguments() []string {
	return *opts.arguments
}

func (opts *options) HasArguments() bool {
	return opts.arguments != nil
}
