package god

import "time"

type intervalOption interface {
	WithInterval(time.Duration) Options
	Interval() time.Duration
	HasInterval() bool
}

var _ intervalOption = &options{}

func (opts *options) WithInterval(interval time.Duration) Options {
	newOpts := opts.copy()
	newOpts.interval = &interval
	return newOpts
}

func (opts *options) Interval() time.Duration {
	return *opts.interval
}
func (opts *options) HasInterval() bool {
	return opts.interval != nil
}
