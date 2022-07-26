package god

type watchingFilesOption interface {
	WithWatchingFiles(...string) Options
	WatchingFiles() []string
	HasWatchingFiles() bool
}

var _ watchingFilesOption = &options{}

func (opts *options) WithWatchingFiles(watchingFiles ...string) Options {
	newOpts := opts.copy()
	newOpts.watchingFiles = &watchingFiles
	return newOpts
}

func (opts *options) WatchingFiles() []string {
	return *opts.watchingFiles
}

func (opts *options) HasWatchingFiles() bool {
	return opts.watchingFiles != nil
}
