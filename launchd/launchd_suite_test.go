package launchd_test

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var bashPath string

var _ = BeforeSuite(func() {
	bashPath = os.Getenv("BASH_PATH")
	if bashPath == "" {
		path, err := exec.LookPath("bash")
		Expect(err).ToNot(HaveOccurred())

		bashPath = path
	}
})

func TestLaunchd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Launchd Suite")
}
