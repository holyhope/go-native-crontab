package god

type scopeOption interface {
	WithScope(Scope) Options
	Scope() Scope
	HasScope() bool
}

var _ scopeOption = &options{}

func (opts *options) WithScope(scope Scope) Options {
	newOpts := opts.copy()
	newOpts.scope = &scope
	return newOpts
}

func (opts *options) Scope() Scope {
	return *opts.scope
}

func (opts *options) HasScope() bool {
	return opts.scope != nil
}

//go:generate stringer -type=Scope -trimprefix=Scope -output options_scope_stringer.go
// ScopeValue is the scope of the unit.
type Scope uint8

const (
	ScopeUser Scope = iota
	ScopeSystem
)
