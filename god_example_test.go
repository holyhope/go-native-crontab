package god_test

import (
	"context"
	"fmt"

	"github.com/holyhope/god"
)

func ExampleNew() {
	ct, err := god.New(context.Background())
	if err != nil {
		panic(err)
	}

	_ = ct

	fmt.Println("New god object created")

	// Output:
	// New god object created
}

func Example() {
	cronTab, _ := god.New(
		context.Background(),
		god.UnitName("com.github.holyhope.test.god_example"),
		god.UnitCommand("/bin/bash", "-c", `echo "Hello, world!"`),
		god.ScopeUser,
	)

	// Install the unit to the system
	_ = cronTab.Install(context.Background())

	fmt.Println("Unit installed")

	// Install the unit to the system
	_ = cronTab.Uninstall(context.Background())

	fmt.Println("Unit uninstalled")

	// Output:
	// Unit installed
	// Unit uninstalled
}
