package god

type userOwnerOption interface {
	WithUserOwner(int) Options
	UserOwner() int
	HasUserOwner() bool
}

var _ userOwnerOption = &options{}

func (opts *options) WithUserOwner(uid int) Options {
	newOpts := opts.copy()
	newOpts.userOwner = &uid
	return newOpts
}

func (opts *options) UserOwner() int {
	return *opts.userOwner
}

func (opts *options) HasUserOwner() bool {
	return opts.userOwner != nil
}
