package launchd_test

import (
	"context"
	"errors"
	"io/fs"
	"syscall"
	"testing"
	"time"

	"github.com/holyhope/god"
	"github.com/holyhope/god/launchd"

	. "github.com/onsi/gomega"
)

func FuzzInstall(f *testing.F) {
	f.Add("test", "/bin/true", uint8(god.ScopeUser), uint8(god.DarwinLimitLoadToSessionBackground), int64(time.Second))

	f.Fuzz(func(t *testing.T, name string, command string, scope uint8, darwinLimit uint8, interval int64) {
		g := NewWithT(t)
		Ω := g.Ω

		scopeValue := god.Scope(scope)
		intervalValue := time.Duration(interval)
		darwinSessionTypeValue := god.DarwinLimitLoadToSessionType(darwinLimit)

		unit, err := launchd.New(
			context.Background(),
			god.Opts().
				WithName(name).
				WithScope(scopeValue).
				WithProgram(command).
				WithInterval(intervalValue).
				WithDarwinLimitLoadToSessionType(darwinSessionTypeValue),
		)

		MatchInvalidOptionError := Or(
			MatchError(&god.InvalidOptionError{
				Key:   "Name",
				Value: name,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   "Program",
				Value: command,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   "Scope",
				Value: scopeValue,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   "Interval",
				Value: intervalValue,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   "DarwinLimitLoadToSessionType",
				Value: darwinSessionTypeValue,
			}),
		)

		if err != nil {
			Ω(err).Should(MatchInvalidOptionError, "New should be successful or fail with InvalidOptionError")

			t.Log("New failed as expected", err)

			return
		}

		Ω(unit).ShouldNot(BeNil())

		t.Logf("Installing %s", name)

		var installationError error

		defer func() {
			t.Logf("Uninstalling %s", name)
			err := unit.Uninstall(context.Background())

			// Accept ENOENT and ExitError due to parallelism: the file can be deleted by another test

			if installationError != nil {
				Ω(err).Should(Or(Succeed(), MatchError(syscall.ENAMETOOLONG), MatchError(syscall.EINVAL), MatchError(syscall.ENOENT), BeAssignableToTypeOf(&launchd.ExecError{})), "Uninstallation should fail with the right error")
			} else if !t.Failed() && !t.Skipped() {
				Ω(err).Should(Or(Succeed(), MatchError(syscall.ENOENT), MatchError(errors.New("Boot-out failed: 5: Input/output error\nTry re-running the command as root for richer errors.\n"))), "Uninstallation should succeed")
			}
		}()

		installationError = unit.Install(context.Background())

		Ω(installationError).Should(Or(Succeed(), MatchError(fs.ErrPermission), MatchError(syscall.EINVAL), MatchError(syscall.EILSEQ), MatchError(syscall.ENAMETOOLONG), MatchError(errors.New("Bootstrap failed: 5: Input/output error\nTry re-running the command as root for richer errors.\n"))), "Installation should succeed or fail with the right error")
	})
}
