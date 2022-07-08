//go:build !darwin && linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build !darwin,linux,!freebsd,!netbsd,!openbsd,!windows,!js

package crontab

import (
	"context"
	"errors"
)

// New creates a new CronTab.
func New(_ context.Context, _ ...FactoryOpts) (CronTab, error) {
	return nil, errors.New("not implemented")
}
