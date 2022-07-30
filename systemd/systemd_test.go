package systemd_test

import (
	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"

	"github.com/holyhope/god/internal/tests"
	"github.com/holyhope/god/systemd"
)

var _ = Describe("Systemd", func() {
	tests.NewSuite(systemd.New)
})
