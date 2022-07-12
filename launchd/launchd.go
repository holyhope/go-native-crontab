//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package launchd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/holyhope/god"
	"howett.net/plist"
)

// New creates a new Unit.
func New(ctx context.Context, opts god.Options) (god.Unit, error) {
	u := &Unit{}

	if opts == nil {
		opts = god.With()
	}

	if err := u.Apply(opts); err != nil {
		return nil, fmt.Errorf("apply options: %w", err)
	}

	return u, nil
}

type Unit struct {
	name             string
	unitType         string
	state            string
	description      string
	program          string
	programArguments []string
	envs             map[string]string
	scope            god.ScopeValue
	interval         time.Duration
}

func (u *Unit) apply(opts god.Options, key god.OptionKey, setter func(god.OptionValue) bool) error {
	value := opts.Get(key)

	if !setter(value) {
		return &god.InvalidOptionError{
			Key:   key,
			Value: value,
		}
	}

	return nil
}

func (u *Unit) Apply(opts god.Options) error {
	if err := u.apply(opts, god.Name, func(value god.OptionValue) bool {
		name, ok := value.(string)
		if !ok {
			return false
		}

		if !IsValidName(name) {
			return false
		}

		u.name = name

		return true
	}); err != nil {
		return err
	}

	if !opts.Has(god.Type) {
		u.unitType = "agent"
	} else if err := u.apply(opts, god.Type, func(value god.OptionValue) bool {
		unitType, ok := value.(string)
		if !ok {
			return false
		}

		if !IsValidType(unitType) {
			return false
		}

		u.unitType = unitType

		return true
	}); err != nil {
		return err
	}

	if err := u.apply(opts, god.Program, func(value god.OptionValue) bool {
		program, ok := value.(string)
		if !ok {
			return false
		}

		if !IsValidProgram(program) {
			return false
		}

		u.program = program

		return true
	}); err != nil {
		return err
	}

	if !opts.Has(god.ProgramArguments) {
		u.programArguments = []string{}
	} else if err := u.apply(opts, god.ProgramArguments, func(value god.OptionValue) bool {
		args, ok := value.([]string)
		if !ok {
			return false
		}

		if !IsValidProgramArguments(args) {
			return false
		}

		u.programArguments = args

		return true
	}); err != nil {
		return err
	}

	if !opts.Has(god.Interval) {
		u.interval = 0
	} else if err := u.apply(opts, god.Interval, func(value god.OptionValue) bool {
		interval, ok := value.(time.Duration)
		if !ok {
			return false
		}

		if !IsValidInterval(interval) {
			return false
		}

		u.interval = interval

		return true
	}); err != nil {
		return err
	}

	return nil
}

func (u *Unit) Name() string {
	return u.name
}

func (u *Unit) Description() string {
	return u.description
}

func (u *Unit) Type() string {
	return u.unitType
}

func (u *Unit) State() string {
	return u.state
}

func (u *Unit) Command() []string {
	return append([]string{u.program}, u.programArguments...)
}

func (u *Unit) Envs() map[string]string {
	return u.envs
}

func (u *Unit) Scope() god.ScopeValue {
	return u.scope
}

func (u *Unit) Install(ctx context.Context) error {
	unitPath, err := writeUnitFile(u)
	if err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	cmd := exec.Command("launchctl", "load", unitPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("load %s: %w", unitPath, err)
	}

	return nil
}

func (u *Unit) Uninstall(ctx context.Context) error {
	unitPath, err := u.path()
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "launchctl", "unload", unitPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unload %s: %w", unitPath, err)
	}

	if err := os.Remove(unitPath); err != nil {
		return fmt.Errorf("remove unit file: %w", err)
	}

	return nil
}

func (u *Unit) path() (string, error) {
	result, err := Path(u.scope, u.name+".plist")
	if err != nil {
		return "", fmt.Errorf("get path: %w", err)
	}

	return result, nil
}

func writeUnitFile(u *Unit) (string, error) {
	unitPath, err := u.path()
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(unitPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return unitPath, &god.InvalidOptionError{
				Key:   god.Scope,
				Value: u.scope,
			}
		}

		pathError := &os.PathError{}
		if errors.As(err, &pathError) {
			return unitPath, &god.InvalidOptionError{
				Key:   god.Name,
				Value: u.name,
			}
		}

		return unitPath, fmt.Errorf("open %s: %w", unitPath, err)
	}

	defer f.Close()

	if err := u.ToPlist(f); err != nil {
		return unitPath, fmt.Errorf("convert to plist: %w", err)
	}

	return unitPath, nil
}

func (u *Unit) ToPlist(writer io.Writer) error {
	agent := make(launchUnit, 6)

	agent.Label(u.name)
	agent.Program(u.program)
	agent.ProgramArguments(u.programArguments...)
	agent.RunAtLoad(true)
	agent.StartInterval(int(u.interval.Seconds()))
	agent.KeepAlive(false)
	agent.EnvironmentVariables(u.envs)

	encoder := plist.NewEncoder(writer)
	if err := encoder.Encode(agent); err != nil {
		return fmt.Errorf("encode agent: %w", err)
	}

	return nil
}
