//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package god

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"howett.net/plist"
)

// New creates a new Unit.
func New(ctx context.Context, opts ...FactoryOpts) (u Unit, err error) {
	u = &unit{}

	for i, opt := range opts {
		if err := opt.Apply(ctx, u); err != nil {
			return nil, fmt.Errorf("apply %d (%+v): %w", i, opt, err)
		}
	}

	return u, nil
}

type unit struct {
	name        string
	unitType    string
	state       string
	description string
	command     []string
	envs        map[string]string
	scope       UnitScope
}

func (u *unit) Name() string {
	return u.name
}

func (u *unit) Description() string {
	return u.description
}

func (u *unit) Type() string {
	return u.unitType
}

func (u *unit) State() string {
	return u.state
}

func (u *unit) Command() []string {
	return u.command
}

func (u *unit) Envs() map[string]string {
	return u.envs
}

func (u *unit) Scope() UnitScope {
	return u.scope
}

type launchUnit struct {
	Label            string   `plist:"Label"`
	Program          string   `plist:"Program"`
	ProgramArguments []string `plist:"ProgramArguments"`
	RunAtLoad        bool     `plist:"RunAtLoad"`
	StartInterval    int      `plist:"StartInterval"`
}

func (u *unit) Install(ctx context.Context) error {
	var missingsError MissingOptionsError

	if u.name == "" {
		missingsError.addMissing("UnitName")
	}

	if u.scope == scopeUnspecified {
		missingsError.addMissing("Scope")
	}

	if !missingsError.IsEmpty() {
		return &missingsError
	}

	path, err := u.scope.Path(u.name + ".plist")
	if err != nil {
		return fmt.Errorf("get path: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	defer f.Close()

	if err := u.ToPlist(f); err != nil {
		return fmt.Errorf("convert to plist: %w", err)
	}

	cmd := exec.Command("launchctl", "load", path)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}

	return nil
}

func (u *unit) ToPlist(writer io.Writer) error {
	agent := &launchUnit{
		Label:            u.name,
		Program:          u.command[0],
		ProgramArguments: u.command[1:],
		RunAtLoad:        true,
		StartInterval:    1,
	}

	encoder := plist.NewEncoder(writer)
	if err := encoder.Encode(agent); err != nil {
		return fmt.Errorf("encode agent: %w", err)
	}

	return nil
}

func (u *unit) Uninstall(ctx context.Context) error {
	path, err := u.scope.Path(u.name + ".plist")
	if err != nil {
		return fmt.Errorf("get path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "launchctl", "unload", path)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unload %s: %w", path, err)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove file %s: %w", path, err)
	}

	return nil
}
