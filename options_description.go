package god

type descriptionOption interface {
	WithDescription(string) Options
	Description() string
	HasDescription() bool
}

var _ descriptionOption = &options{}

func (opts *options) WithDescription(description string) Options {
	newOpts := opts.copy()
	newOpts.description = &description
	return newOpts
}

func (opts *options) Description() string {
	return *opts.description
}
func (opts *options) HasDescription() bool {
	return opts.description != nil
}
