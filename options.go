package god

import (
	"context"
	"fmt"
	"time"
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

// Name is the name of the unit.
type Name string

// Type is the type of the unit.
type Type string

// State is the state of the unit.
type State string

// Description is the description of the unit.
type Description string

// Envs is the environment variables of the unit.
type Envs map[string]string

// Interval is the interval of the unit.
type Interval time.Duration

// UnitScope is the scope of the unit.
//go:generate stringer -type=Scope -linecomment
type Scope uint8

const (
	scopeUnspecified Scope = iota // unspecified
	ScopeUser                     // user
	ScopeSystem                   // system
)

type UnkownScopeError Scope

func (err UnkownScopeError) Error() string {
	return fmt.Sprintf("unknown scope %s", Scope(err).String())
}
