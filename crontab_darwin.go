//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package crontab

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"howett.net/plist"
)

// New creates a new CronTab.
func New(ctx context.Context, opts ...FactoryOpts) (ct CronTab, err error) {
	ct = &cronTab{
		Entries:  []*cronEntry{},
		Scope:    darwinNoScope,
		FileName: "",
	}

	for i, opt := range opts {
		ct, err = opt.Apply(ctx, ct)
		if err != nil {
			return nil, fmt.Errorf("apply %d (%+v): %w", i, opt, err)
		}
	}

	return ct, nil
}

type cronTab struct {
	Entries  []*cronEntry
	Scope    DarwinOptsScope
	FileName string
}

type cronEntry struct {
	Label            string   `plist:"Label"`
	Program          string   `plist:"Program"`
	ProgramArguments []string `plist:"ProgramArguments"`
	RunAtLoad        bool     `plist:"RunAtLoad"`
	StartInterval    int      `plist:"StartInterval"`
}

func (ct *cronTab) Add(ctx context.Context, interval time.Duration, program string, args ...string) error {
	ct.Entries = append(ct.Entries, &cronEntry{
		Label:            "crontab",
		Program:          program,
		ProgramArguments: args,
		RunAtLoad:        true,
		StartInterval:    int(interval.Seconds()),
	})

	return nil
}

var optIndex interface{}

func optionIndex(ctx context.Context) int {
	index := ctx.Value(&optIndex)
	if index == nil {
		return 0
	}

	return index.(int) // nolint: forcetypeassert
}

func incrementOptionIndex(ctx context.Context) context.Context {
	return context.WithValue(ctx, &optIndex, optionIndex(ctx)+1)
}

func (ct *cronTab) Install(ctx context.Context, opts ...InstallOpts) (InstalledCronTab, error) {
	if len(opts) > 0 {
		computedCronTab, err := opts[0].Apply(ctx, ct)
		if err != nil {
			return nil, fmt.Errorf("apply %d (%+v): %w", optionIndex(ctx), opts[0], err)
		}

		return computedCronTab.Install(incrementOptionIndex(ctx), opts[1:]...) // nolint: wrapcheck
	}

	destination, err := ct.destination()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", destination, err)
	}

	defer f.Close()

	if err := ct.ToPlist(f); err != nil {
		return nil, fmt.Errorf("convert to plist: %w", err)
	}

	return &installedCronTab{ct, destination}, nil
}

type UnkownScopeError DarwinOptsScope

func (err UnkownScopeError) Error() string {
	return fmt.Sprintf("unknown scope %s", DarwinOptsScope(err).String())
}

func (ct *cronTab) destination() (string, error) {
	var destination string

	if ct.FileName == "" {
		return "", &MissingOptionsError{
			Name: "FileName",
		}
	}

	switch ct.Scope {
	case UserScope:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home: %w", err)
		}

		destination = path.Join(home, "Library/LaunchAgents", ct.FileName+".plist")
	case SystemScope:
		destination = path.Join("/Library/LaunchAgents", ct.FileName+".plist")
	case DarwinDaemonSystemScope:
		destination = path.Join("/Library/LaunchDaemons/crontab.plist", ct.FileName+".plist")
	case darwinNoScope:
		return "", &MissingOptionsError{
			Name: "Scope",
		}
	default:
		return "", UnkownScopeError(ct.Scope)
	}

	return destination, nil
}

func (ct *cronTab) ToPlist(writer io.Writer) error {
	agent := &cronEntry{
		Label:            "crontab",
		Program:          "/bin/sh",
		ProgramArguments: []string{"-c", "crontab"},
		RunAtLoad:        true,
		StartInterval:    1,
	}

	encoder := plist.NewEncoder(writer)
	if err := encoder.Encode(agent); err != nil {
		return fmt.Errorf("encode agent: %w", err)
	}

	return nil
}

type installedCronTab struct {
	*cronTab
	path string
}

func (ct *installedCronTab) Path() string {
	return ct.path
}

func (ct *installedCronTab) Uninstall(ctx context.Context) error {
	destination, err := ct.destination()
	if err != nil {
		return err
	}

	if err := os.Remove(destination); err != nil {
		return fmt.Errorf("remove file %s: %w", destination, err)
	}

	return nil
}
