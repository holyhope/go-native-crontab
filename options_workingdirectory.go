package god

type workingDirectoryOption interface {
	WithWorkingDirectory(string) Options
	WorkingDirectory() string
	HasWorkingDirectory() bool
}

var _ programOption = &options{}

func (opts *options) WithWorkingDirectory(workingDirectory string) Options {
	newOpts := opts.copy()
	newOpts.workingDirectory = &workingDirectory
	return newOpts
}

func (opts *options) WorkingDirectory() string {
	return *opts.workingDirectory
}

func (opts *options) HasWorkingDirectory() bool {
	return opts.workingDirectory != nil
}
