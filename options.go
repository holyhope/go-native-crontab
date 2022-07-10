package god

import (
	"context"
	"fmt"
)

type FactoryOptsFn func(context.Context, Unit) error

func (f FactoryOptsFn) Apply(ctx context.Context, u Unit) error {
	if f == nil {
		return nil
	}

	return f(ctx, u)
}

func NoOp(ctx context.Context, ct Unit) (Unit, error) {
	return ct, nil
}

type MissingOptionsError struct {
	Missings []string
}

func (err *MissingOptionsError) Error() string {
	return fmt.Sprintf("missing option %v", err.Missings)
}

func (err *MissingOptionsError) addMissing(name string) {
	err.Missings = append(err.Missings, name)
}

func (err *MissingOptionsError) IsEmpty() bool {
	return len(err.Missings) == 0
}

// UnitName is the name of the unit.
type UnitName string

// UnitType is the type of the unit.
type UnitType string

// UnitState is the state of the unit.
type UnitState string

// UnitDescription is the description of the unit.
type UnitDescription string

// UnitEnvs is the environment variables of the unit.
type UnitEnvs map[string]string

// UnitScope is the scope of the unit.
// go:generate stringer -type=UnitScope -linecomment
type UnitScope uint8

const (
	scopeUnspecified UnitScope = iota // unspecified
	ScopeUser                         // user
	ScopeSystem                       // system
)

type UnkownScopeError UnitScope

func (err UnkownScopeError) Error() string {
	return fmt.Sprintf("unknown scope %s", UnitScope(err).String())
}
