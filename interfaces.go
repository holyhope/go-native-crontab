package god

import (
	"context"
)

// Unit represents single service managed by either systemd or launchd.
type Unit interface {
	// Install the unit to the system.
	Install(ctx context.Context) error
	// Uninstall the unit from the system.
	Uninstall(ctx context.Context) error
	// Status the status of the unit.
	Status(ctx context.Context) (UnitStatus, error)
}

type UnitStatus interface {
	// Exists returns true if the unit is on the filesystem.
	Exists(ctx context.Context) bool
	// IsLoaded returns true if the unit is loaded.
	IsLoaded(ctx context.Context) bool
}
