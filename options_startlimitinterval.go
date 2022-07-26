package god

import "time"

type startLimitIntervalOption interface {
	WithStartLimitInterval(startLimitInterval time.Duration) Options
	StartLimitInterval() time.Duration
	HasStartLimitInterval() bool
}

var _ startLimitIntervalOption = &options{}

func (opts *options) WithStartLimitInterval(startLimitInterval time.Duration) Options {
	newOpts := opts.copy()
	newOpts.startLimitInterval = &startLimitInterval
	return newOpts
}

func (opts *options) StartLimitInterval() time.Duration {
	return *opts.startLimitInterval
}

func (opts *options) HasStartLimitInterval() bool {
	return opts.startLimitInterval != nil
}
