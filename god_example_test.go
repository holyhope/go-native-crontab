package god_test

import (
	"context"
	"fmt"

	"github.com/holyhope/god"
	_ "github.com/holyhope/god/launchd"
)

func ExampleNew() {
	god.New(context.Background(), god.With(god.Name, "test"))

	fmt.Printf("New god object created")

	// Output:
	// New god object created
}

func Example() {
	unit, err := god.New(
		context.Background(),
		god.With(
			god.Name, "com.github.holyhope.test.god_example",
			god.Program, "/bin/bash",
			god.ProgramArguments, []string{"-c", `echo "Hello, world!"`},
			god.Scope, god.ScopeUser,
		),
	)
	if err != nil {
		panic(err)
	}

	// Install the unit to the system
	_ = unit.Install(context.Background())

	fmt.Println("Unit installed")

	// Install the unit to the system
	_ = unit.Uninstall(context.Background())

	fmt.Println("Unit uninstalled")

	// Output:
	// Unit installed
	// Unit uninstalled
}
