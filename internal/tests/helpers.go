package tests

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/holyhope/god"
	"github.com/iancoleman/strcase"
)

func BashPath() string {
	bashPath := os.Getenv("BASH_PATH")
	if bashPath != "" {
		return bashPath
	}

	path, err := exec.LookPath("bash")
	Expect(err).ToNot(HaveOccurred())

	return path
}

func NewSuite(factory func(ctx context.Context, opts god.Options) (god.Unit, error)) {
	Describe("New unit", Offset(1), func() {
		var opts god.Options

		BeforeEach(func() {
			opts = nil
		})

		Context("With no options", func() {
			It("should return an error", func() {
				_, err := factory(context.Background(), opts)
				Expect(err).To(Or(
					MatchError(&god.MissingOptionError{
						Key: "Name",
					}),
					MatchError(&god.MissingOptionError{
						Key: "Program",
					}),
					MatchError(&god.MissingOptionError{
						Key: "Scope",
					}),
				))
			})
		})
	})

	Describe("Install", Offset(1), func() {
		var unit god.Unit

		AfterEach(func() {
			Ω(unit.Uninstall(context.Background())).Should(Succeed())
		})

		Context("With all options", func() {
			BeforeEach(func() {
				name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

				var err error
				unit, err = factory(context.Background(), god.Opts().
					WithName(name).
					WithScope(god.ScopeUser).
					WithProgram(BashPath()).
					WithArguments("-c", `echo 'Hello, world!'`).
					WithUserOwner(os.Getuid()).
					WithDarwinLimitLoadToSessionType(god.DarwinLimitLoadToSessionBackground),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(unit).ToNot(BeNil())
			})

			It("Should work", func() {
				Ω(unit.Install(context.Background())).Should(Succeed())
			})

			It("Can be installed multiple times", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Install(context.Background())).Should(Succeed())
				Ω(unit.Install(context.Background())).Should(Succeed())
			})
		})
	})

	Describe("Uninstall", Offset(1), func() {
		var unit god.Unit

		BeforeEach(func() {
			name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

			var err error
			unit, err = factory(context.Background(), god.Opts().
				WithName(name).
				WithScope(god.ScopeUser).
				WithProgram(BashPath()).
				WithArguments("-c", `echo 'Hello, world!'`).
				WithUserOwner(os.Getuid()).
				WithDarwinLimitLoadToSessionType(god.DarwinLimitLoadToSessionBackground),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(unit).ToNot(BeNil())
		})

		Context("without installation", func() {
			It("should work", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Uninstall(context.Background())).Should(Succeed())
			})
		})

		Context("After installation", func() {
			BeforeEach(func() {
				Ω(unit.Install(context.Background())).Should(Succeed())
			})

			It("should work", func() {
				Ω(unit.Uninstall(context.Background())).Should(Succeed())
			})
		})
	})
}
