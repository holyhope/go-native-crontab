package crontab_test

import (
	"context"

	crontab "github.com/holyhope/go-native-crontab"
)

func ExampleFileName() {
	ct, _ := crontab.New(context.Background())

	ict, err := ct.Install(
		context.Background(),
		crontab.FileName("com.github.holyhope.crontab_example"),
	)

	if err == nil {
		// Cleanup filesystem
		ict.Uninstall(context.Background())
	}
}

func ExampleUserScope() {
	ct, _ := crontab.New(
		context.Background(),
		crontab.UserScope,
	)

	ict, err := ct.Install(
		context.Background(),
	)
	if err == nil {
		// Cleanup filesystem
		ict.Uninstall(context.Background())
	}
}

func ExampleSystemScope() {
	ct, _ := crontab.New(
		context.Background(),
		crontab.SystemScope,
	)

	ict, err := ct.Install(
		context.Background(),
	)
	if err == nil {
		// Cleanup filesystem
		ict.Uninstall(context.Background())
	}
}
