//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package launchd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/holyhope/god"
	"howett.net/plist"
)

// New creates a new Unit.
func New(ctx context.Context, opts god.Options) (god.Unit, error) {
	if opts == nil {
		opts = god.Opts()
	}

	name, err := Name(opts)
	if err != nil {
		return nil, err
	}

	scope, err := Scope(opts)
	if err != nil {
		return nil, err
	}

	unitP, err := Path(scope, name+".plist")
	if err != nil {
		return nil, err
	}

	launchU := make(launchUnit, 2)

	launchU.Label(name)

	if !opts.HasProgram() {
		return nil, god.NewMissingOptionError("Program")
	}
	launchU.Program(opts.Program())
	if opts.HasDarwinLimitLoadToSessionType() {
		limitLoadToSessionType, err := LimitLoadToSessionType(opts.DarwinLimitLoadToSessionType())
		if err != nil {
			return nil, err
		}

		launchU.LimitLoadSessionType(limitLoadToSessionType)
	}

	launchU.KeepAlive(false)
	if opts.HasArguments() {
		launchU.ProgramArguments(opts.Arguments()...)
	}
	if opts.HasRunAtLoad() {
		launchU.RunAtLoad(opts.RunAtLoad())
	}
	if opts.HasWorkingDirectory() {
		launchU.WorkingDirectory(opts.WorkingDirectory())
	}
	if opts.HasStandardOutput() {
		launchU.StandardOutPath(opts.StandardOutput())
	}
	if opts.HasErrorOutput() {
		launchU.StandardErrorPath(opts.ErrorOutput())
	}
	if opts.HasInterval() {
		launchU.StartInterval(int(opts.Interval().Seconds()))
	}
	if opts.HasEnvironmentVariables() {
		launchU.EnvironmentVariables(opts.EnvironmentVariables())
	}
	if opts.HasEnvironmentVariables() {
		launchU.EnvironmentVariables(opts.EnvironmentVariables())
	}
	if opts.HasUserOwner() {
		username, err := Username(opts)
		if err != nil {
			return nil, err
		}

		launchU.Username(username)
	}

	launchctlPath, err := exec.LookPath("launchctl")
	if err != nil {
		return nil, fmt.Errorf("launchctl not found: %w", err)
	}

	domain, err := Domain(scope)
	if err != nil {
		return nil, err
	}

	return &Unit{
		unitSpec:      launchU,
		domain:        domain,
		unitPath:      unitP,
		launchctlPath: launchctlPath,
	}, nil
}

type Unit struct {
	unitSpec      launchUnit
	domain        string
	unitPath      string
	launchctlPath string
}

func (u *Unit) Install(ctx context.Context) error {
	if err := u.writeUnitFile(u.unitPath); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	if err := u.exec(ctx, u.launchctlPath, "bootstrap", u.domain, u.unitPath); err != nil {
		if err, ok := err.(*ExecError); ok {
			if err.MatchLaunchdReason("service already bootstrapped") {
				return nil
			}
		}

		return err
	}

	return nil
}

func (u *Unit) exec(ctx context.Context, command string, args ...string) error {
	var stderrReader, stdoutReader io.Reader

	cmd := exec.CommandContext(ctx, command, args...)

	r, err := cmd.StderrPipe()
	if err != nil {
		stderrReader = bytes.NewBufferString(fmt.Sprintf(`failed to create stderr pipe: %s`, err))
	} else {
		defer r.Close()

		stderrReader = r
	}

	r, err = cmd.StdoutPipe()
	if err != nil {
		stdoutReader = bytes.NewBufferString(fmt.Sprintf(`failed to create stdout pipe: %s`, err))
	} else {
		defer r.Close()

		stdoutReader = r
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start launchctl: %w", err)
	}

	stderr, err := io.ReadAll(stderrReader)
	if err != nil {
		stderr = []byte(fmt.Sprintf(`failed to read stderr: %s`, err))
	}

	stdout, err := io.ReadAll(stdoutReader)
	if err != nil {
		stdout = []byte(fmt.Sprintf(`failed to read stdout: %s`, err))
	}

	if err := cmd.Wait(); err != nil {
		return &ExecError{
			UnderlyingError: err,
			Command:         cmd.Args,
			Stderr:          string(stderr),
			Stdout:          string(stdout),
		}
	}

	return nil
}

func (u *Unit) Uninstall(ctx context.Context) error {
	err := u.exec(ctx, u.launchctlPath, "bootout", u.domain, u.unitPath)
	if err != nil {
		if err, ok := err.(*ExecError); ok {
			if err.MatchLaunchdReason("No such file or directory") {
				return nil
			}
		}

		return err
	}

	if err := os.Remove(u.unitPath); err != nil {
		return fmt.Errorf("remove unit file: %w", err)
	}

	return nil
}

func (u *Unit) writeUnitFile(unitPath string) error {
	f, err := os.OpenFile(unitPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", unitPath, err)
	}

	defer f.Close()

	if err := u.ToPlist(f); err != nil {
		return fmt.Errorf("convert to plist: %w", err)
	}

	return nil
}

func (u *Unit) ToPlist(writer io.Writer) error {
	encoder := plist.NewEncoder(writer)
	if err := encoder.Encode(u.unitSpec); err != nil {
		return fmt.Errorf("encode unit: %w", err)
	}

	return nil
}
