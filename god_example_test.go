package god_test

import (
	"context"
	"fmt"

	"github.com/holyhope/god"
	_ "github.com/holyhope/god/launchd"
)

func ExampleNew() {
	god.New(context.Background(), god.Opts().WithName("test"))

	fmt.Printf("New god object created")

	// Output:
	// New god object created
}

func Example() {
	unit, err := god.New(
		context.Background(),
		god.Opts().
			WithName("com.github.holyhope.test.god_example").
			WithProgram("/bin/bash").
			WithArguments("-c", `echo "Hello, world!"`).
			WithScope(god.ScopeUser),
	)
	if err != nil {
		panic(err)
	}

	// Install the unit to the system
	_ = unit.Create(context.Background())

	fmt.Println("Unit installed")

	// Enable the unit
	_ = unit.Enable(context.Background())

	fmt.Println("Unit enabled")

	// Install the unit to the system
	_ = unit.Delete(context.Background())

	fmt.Println("Unit uninstalled")

	// Output:
	// Unit installed
	// Unit enabled
	// Unit uninstalled
}
