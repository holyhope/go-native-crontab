//go:build !darwin && linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build !darwin,linux,!freebsd,!netbsd,!openbsd,!windows,!js

package god

import (
	"context"
	"errors"
)

// New creates a new CronTab.
func New(_ context.Context, _ ...FactoryOpts) (Unit, error) {
	return nil, errors.New("not implemented")
}
