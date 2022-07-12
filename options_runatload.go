package god

type runAtLoadOption interface {
	WithRunAtLoad(bool) Options
	RunAtLoad() bool
	HasRunAtLoad() bool
}

var _ runAtLoadOption = &options{}

func (opts *options) WithRunAtLoad(runAtLoad bool) Options {
	newOpts := opts.copy()
	newOpts.runAtLoad = &runAtLoad
	return newOpts
}

func (opts *options) RunAtLoad() bool {
	return *opts.runAtLoad
}
func (opts *options) HasRunAtLoad() bool {
	return opts.runAtLoad != nil
}
