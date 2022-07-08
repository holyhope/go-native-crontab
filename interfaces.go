package crontab

import (
	"context"
	"time"
)

// CronTab is a crontab.
type CronTab interface {
	// Add adds a new cron entry to the crontab.
	Add(context.Context, time.Duration, string, ...string) error
	// Install installs the crontab to the current system.
	Install(context.Context, ...InstallOpts) (InstalledCronTab, error)
}

// InstalledCronTab is a crontab that has been installed to the current system.
type InstalledCronTab interface {
	// Remove removes the crontab from the current system.
	Uninstall(context.Context) error
	Path() string
}

// FactoryOpts is an option for installing a crontab.
type FactoryOpts interface {
	// Apply applies the install options to the crontab.
	Apply(context.Context, CronTab) (CronTab, error)
}

// InstallOpts is an option for installing a crontab.
type InstallOpts interface {
	// Apply applies the install options to the crontab.
	Apply(context.Context, CronTab) (CronTab, error)
}
