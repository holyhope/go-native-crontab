package crontab_test

import (
	"context"
	"fmt"
	"time"

	crontab "github.com/holyhope/go-native-crontab"
)

func ExampleNew() {
	ct, err := crontab.New(context.Background())
	if err != nil {
		panic(err)
	}

	_ = ct

	fmt.Println("New crontab object created")

	// Output:
	// New crontab object created
}

func Example() {
	cronTab, _ := crontab.New(
		context.Background(),
		crontab.UserScope,
	)

	// Run the bach command every minute
	cronTab.Add(context.Background(), time.Minute, "bash", "-c", "echo 'Hello world!'")

	// Install the crontab to the system
	ict, err := cronTab.Install(
		context.Background(),
		crontab.FileName("com.github.holyhope.crontab_example"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Crontab installed")

	if err := ict.Uninstall(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println("Crontab uninstalled")

	// Output:
	// Crontab installed
	// Crontab uninstalled
}
