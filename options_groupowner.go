package god

type groupOwnerOption interface {
	WithGroupOwner(int) Options
	GroupOwner() int
	HasGroupOwner() bool
}

var _ groupOwnerOption = &options{}

func (opts *options) WithGroupOwner(gid int) Options {
	newOpts := opts.copy()
	newOpts.groupOwner = &gid
	return newOpts
}

func (opts *options) GroupOwner() int {
	return *opts.groupOwner
}

func (opts *options) HasGroupOwner() bool {
	return opts.groupOwner != nil
}
