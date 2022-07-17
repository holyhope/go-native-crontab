package god_test

import (
	"context"

	"github.com/holyhope/god"
	_ "github.com/holyhope/god/launchd"
)

func ExampleOptions() {
	u, _ := god.New(
		context.Background(),
		god.Opts().
			WithName("my-unit").
			WithScope(god.ScopeUser),
	)

	_ = u.Create(context.Background())
	_ = u.Enable(context.Background())

	// Cleanup filesystem
	_ = u.Delete(context.Background())
}
