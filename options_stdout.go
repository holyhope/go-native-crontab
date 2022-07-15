package god

type standardOutputOption interface {
	WithStandardOutput(string) Options
	StandardOutput() string
	HasStandardOutput() bool
}

var _ standardOutputOption = &options{}

func (opts *options) WithStandardOutput(standardOutput string) Options {
	newOpts := opts.copy()
	newOpts.standardOutput = &standardOutput
	return newOpts
}

func (opts *options) StandardOutput() string {
	return *opts.standardOutput
}

func (opts *options) HasStandardOutput() bool {
	return opts.standardOutput != nil
}
