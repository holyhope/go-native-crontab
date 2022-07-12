package god

type Options interface {
	And(...interface{}) Options
	Has(OptionKey) bool
	Get(OptionKey) OptionValue
	Keys() []OptionKey
}

func With(keyValues ...interface{}) Options {
	return (&options{}).And(keyValues...)
}

type options struct {
	Values map[OptionKey]OptionValue
}

func (opts *options) And(keyValues ...interface{}) Options {
	if opts == nil {
		return With(keyValues)
	}

	values := make(map[OptionKey]OptionValue, len(opts.Values)+1)
	for k, v := range opts.Values {
		values[k] = v
	}

	for i := 0; i < len(keyValues); i += 2 {
		values[keyValues[i].(OptionKey)] = keyValues[i+1]
	}

	return &options{Values: values}
}

func (opts *options) Has(key OptionKey) bool {
	_, ok := opts.Values[key]

	return ok
}

func (opts *options) Keys() []OptionKey {
	keys := make([]OptionKey, 0, len(opts.Values))

	for key := range opts.Values {
		keys = append(keys, key)
	}

	return keys
}

func (opts *options) Get(key OptionKey) OptionValue {
	return opts.Values[key]
}

//go:generate stringer -type=OptionKey -linecomment
type OptionKey uint8

const (
	Name OptionKey = iota
	Type
	State
	Description
	Program
	ProgramArguments
	Envs
	Scope
	Interval
)

type OptionValue interface{}

//go:generate stringer -type=ScopeValue -linecomment
// ScopeValue is the scope of the unit.
type ScopeValue uint8

const (
	ScopeUser   ScopeValue = iota // user
	ScopeSystem                   // system
)
