package god_test

import (
	"context"

	god "github.com/holyhope/god"
)

func ExampleName() {
	u, _ := god.New(
		context.Background(),
		god.Name("com.github.holyhope.god_example"),
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}

func ExampleScopeUser() {
	u, _ := god.New(
		context.Background(),
		god.ScopeUser,
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}

func ExampleScopeSystem() {
	u, _ := god.New(
		context.Background(),
		god.ScopeSystem,
	)

	_ = u.Install(context.Background())

	// Cleanup filesystem
	_ = u.Uninstall(context.Background())
}
