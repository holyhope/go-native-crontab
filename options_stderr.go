package god

type errorOutputOption interface {
	WithErrorOutput(string) Options
	ErrorOutput() string
	HasErrorOutput() bool
}

var _ errorOutputOption = &options{}

func (opts *options) WithErrorOutput(errorOutput string) Options {
	newOpts := opts.copy()
	newOpts.errorOutput = &errorOutput
	return newOpts
}

func (opts *options) ErrorOutput() string {
	return *opts.errorOutput
}

func (opts *options) HasErrorOutput() bool {
	return opts.errorOutput != nil
}
