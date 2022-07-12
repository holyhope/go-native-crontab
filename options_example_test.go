package god_test

import (
	"context"

	"github.com/holyhope/god"
	_ "github.com/holyhope/god/launchd"
)

func ExampleName() {
	u, _ := god.New(
		context.Background(),
		god.With(god.Name, "my-unit"),
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}

func ExampleScopeUser() {
	u, _ := god.New(
		context.Background(),
		god.With(god.Scope, god.ScopeUser),
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}

func ExampleScopeSystem() {
	u, _ := god.New(
		context.Background(),
		god.With(god.Scope, god.ScopeSystem),
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}
