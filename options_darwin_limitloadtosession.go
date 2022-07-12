package god

type darwinLimitLoadToSessionTypeOption interface {
	WithDarwinLimitLoadToSessionType(DarwinLimitLoadToSessionType) Options
	DarwinLimitLoadToSessionType() DarwinLimitLoadToSessionType
	HasDarwinLimitLoadToSessionType() bool
}

var _ darwinLimitLoadToSessionTypeOption = &options{}

func (opts *options) WithDarwinLimitLoadToSessionType(sessionType DarwinLimitLoadToSessionType) Options {
	newOpts := opts.copy()
	newOpts.darwin.limitLoadToSessionTypeOption = &sessionType
	return newOpts
}

func (opts *options) DarwinLimitLoadToSessionType() DarwinLimitLoadToSessionType {
	return *opts.darwin.limitLoadToSessionTypeOption
}

func (opts *options) HasDarwinLimitLoadToSessionType() bool {
	return opts.darwin.limitLoadToSessionTypeOption != nil
}

//go:generate stringer -type=DarwinLimitLoadToSessionType -trimprefix=DarwinLimitLoadToSession -output options_darwin_limitloadtosession_stringer.go
// https://developer.apple.com/library/archive/technotes/tn2083/_index.html#//apple_ref/doc/uid/DTS10003794-CH1-SUBSUBSECTION5
type DarwinLimitLoadToSessionType uint8

const (
	DarwinLimitLoadToSessionAqua DarwinLimitLoadToSessionType = iota
	DarwinLimitLoadToSessionStandardIO
	DarwinLimitLoadToSessionBackground
	DarwinLimitLoadToSessionLoginWindow
)
