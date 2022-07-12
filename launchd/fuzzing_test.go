package launchd_test

import (
	"context"
	"testing"
	"time"

	"github.com/holyhope/god"
	"github.com/holyhope/god/launchd"

	. "github.com/onsi/gomega"
)

func FuzzInstall(f *testing.F) {
	f.Add("test", "/bin/true", uint8(god.ScopeUser), int64(time.Second))

	f.Fuzz(func(t *testing.T, name string, command string, scope uint8, interval int64) {
		g := NewWithT(t)
		Ω := g.Ω

		unit, err := launchd.New(
			context.Background(),
			god.With(god.Name, name).
				And(god.Scope, scope).
				And(god.Program, command).
				And(god.Interval, interval),
		)

		MatchInvalidOptionError := Or(
			MatchError(&god.InvalidOptionError{
				Key:   god.Name,
				Value: name,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   god.Program,
				Value: command,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   god.Scope,
				Value: scope,
			}),
			MatchError(&god.InvalidOptionError{
				Key:   god.Interval,
				Value: interval,
			}),
		)

		if err != nil {
			Ω(err).Should(MatchInvalidOptionError, "New should be successful or fail with InvalidOptionError")

			t.Log("New failed as expected", err)

			return
		}

		Ω(unit).ShouldNot(BeNil())

		t.Logf("Installing %s", name)

		defer func() {
			t.Logf("Uninstalling %s", name)
			Ω(unit.Uninstall(context.Background())).Should(Succeed(), "Uninstallation should succeed")
		}()

		Ω(unit.Install(context.Background())).Should(Succeed(), "Installation should succeed")
	})
}
