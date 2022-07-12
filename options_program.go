package god

type programOption interface {
	WithProgram(string) Options
	Program() string
	HasProgram() bool
}

var _ programOption = &options{}

func (opts *options) WithProgram(program string) Options {
	newOpts := opts.copy()
	newOpts.program = &program
	return newOpts
}

func (opts *options) Program() string {
	return *opts.program
}

func (opts *options) HasProgram() bool {
	return opts.program != nil
}
