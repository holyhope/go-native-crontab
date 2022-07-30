package systemd

import (
	"bytes"
	"context"
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

	launchU := make(unit, 2)

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

	// Prepend program name to arguments.
	// See execvp(3) which provide an array of pointers to null-terminated strings that represent the argument list available to the new program.
	// The first argument, by convention, should point to the file name associated with the file being executed.
	if opts.HasArguments() {
		launchU.ExecStart(opts.Program(), opts.Arguments()...)
	} else {
		launchU.ExecStart(opts.Program())
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
		n := time.Now()
		launchU.OnCalendar(fmt.Sprintf("%d %d %d %d %d %d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second()))
	}
	if opts.HasEnvironmentVariables() {
		launchU.EnvironmentVariables(opts.EnvironmentVariables())
	}
	if opts.HasEnvironmentVariables() {
		launchU.EnvironmentVariables(opts.EnvironmentVariables())
	}
	if opts.HasStartLimitInterval() {
		launchU.ThrottleInterval(int(opts.StartLimitInterval().Seconds()))
	}
	if opts.HasWatchingFiles() {
		launchU.WatchPaths(opts.WatchingFiles()...)
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
		unitName:      name,
		domain:        domain,
		unitPath:      unitP,
		launchctlPath: launchctlPath,
	}, nil
}

type Unit struct {
	unitSpec      launchUnit
	unitName      string
	domain        string
	unitPath      string
	launchctlPath string
}

func (u *Unit) exec(ctx context.Context, command string, args ...string) ([]byte, error) {
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
		return nil, fmt.Errorf("start launchctl: %w", err)
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
		return nil, &ExecError{
			UnderlyingError: err,
			Command:         cmd.Args,
			Stderr:          string(stderr),
		}
	}

	return stdout, nil
}

func (u *Unit) writeUnitFile(unitPath string) error {
	f, err := os.OpenFile(unitPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", unitPath, err)
	}

	defer f.Close()

	if err := u.EncodeUnit(f); err != nil {
		return fmt.Errorf("encode unit: %w", err)
	}

	return nil
}

func (u *Unit) EncodeUnit(writer io.Writer) error {
	encoder := plist.NewEncoder(writer)
	if err := encoder.Encode(u.unitSpec); err != nil {
		return fmt.Errorf("encode unit: %w", err)
	}

	return nil
}
