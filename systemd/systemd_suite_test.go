package systemd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLaunchd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Systemd Suite")
}
