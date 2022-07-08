package crontab

import (
	"context"
	"fmt"
)

type InstallOptsFn func(context.Context, CronTab) (CronTab, error)

func (f InstallOptsFn) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	if f == nil {
		return ct, nil
	}

	return f(ctx, ct)
}

func NoOp(ctx context.Context, ct CronTab) (CronTab, error) {
	return ct, nil
}

type MissingOptionsError struct {
	Name string
}

func (err *MissingOptionsError) Error() string {
	return fmt.Sprintf("missing option %s", err.Name)
}

//go:generate stringer -type=DarwinOptsScope -linecomment
type DarwinOptsScope uint32

const (
	darwinNoScope           DarwinOptsScope = iota // Unspecified
	UserScope                                      // User
	SystemScope                                    // System
	DarwinDaemonSystemScope                        // Daemon system
)

type FileName string
