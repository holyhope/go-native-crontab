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
}
