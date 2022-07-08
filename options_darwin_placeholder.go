//go:build !darwin && linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build !darwin,linux,!freebsd,!netbsd,!openbsd,!windows,!js

package crontab

import (
	"context"
)

func (s DarwinOptsScope) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	return ct, nil
}

func (f FileName) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	return ct, nil
}
