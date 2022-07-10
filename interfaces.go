package god

import "context"

// Unit represents single service managed by either systemd or launchd.
type Unit interface {
	// Name is the name of the unit.
	Name() string
	// Description is the description of the unit.
	Description() string
	// Type is the type of the unit.
	Type() string
	// State is the state of the unit.
	State() string
	// Command is the command of the unit.
	Command() []string
	// Envs is the environment variables of the unit.
	Envs() map[string]string
	// Scope is the scope of the unit.
	Scope() Scope
	// Install the unit to the system.
	Install(ctx context.Context) error
	// Uninstall the unit from the system.
	Uninstall(ctx context.Context) error
}

// FactoryOpts is an option for installing a crontab.
type FactoryOpts interface {
	// Apply applies the install options to the crontab.
	Apply(context.Context, Unit) error
}
