package god

type nameOption interface {
	WithName(string) Options
	Name() string
	HasName() bool
}

var _ nameOption = &options{}

func (opts *options) WithName(name string) Options {
	newOpts := opts.copy()
	newOpts.name = &name
	return newOpts
}

func (opts *options) Name() string {
	return *opts.name
}

func (opts *options) HasName() bool {
	return opts.name != nil
}
