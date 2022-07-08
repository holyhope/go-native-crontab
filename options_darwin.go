//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package crontab

import (
	"context"
)

func (s DarwinOptsScope) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	if ct, ok := ct.(*cronTab); ok {
		return &cronTab{
			Entries:  ct.Entries,
			FileName: ct.FileName,
			Scope:    s,
		}, nil
	}

	return ct, nil
}

func (f FileName) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	if ct, ok := ct.(*cronTab); ok {
		return &cronTab{
			Entries:  ct.Entries,
			Scope:    ct.Scope,
			FileName: string(f),
		}, nil
	}

	return ct, nil
}
