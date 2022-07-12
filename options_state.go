package god

type stateOption interface {
	WithState(State) Options
	State() State
	HasState() bool
}

var _ stateOption = &options{}

func (opts *options) WithState(state State) Options {
	newOpts := opts.copy()
	newOpts.state = &state
	return newOpts
}

func (opts *options) State() State {
	return *opts.state
}

func (opts *options) HasState() bool {
	return opts.state != nil
}

//go:generate stringer -type=State -trimprefix=State -output options_state_stringer.go
// ScopeValue is the scope of the unit.
type State uint8

const (
	StateEnable State = iota
	StateDisable
)
