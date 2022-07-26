package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"time"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

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

func CurlPath() string {
	bashPath := os.Getenv("CURL_PATH")
	if bashPath != "" {
		return bashPath
	}

	path, err := exec.LookPath("curl")
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

	Describe("Create", Offset(1), func() {
		var unit god.Unit

		AfterEach(func() {
			Ω(unit.Delete(context.Background())).Should(Succeed())
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
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			It("Can be created multiple times", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Create(context.Background())).Should(Succeed())
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			Context("A previously deleted unit", func() {
				BeforeEach(func() {
					Ω(unit.Create(context.Background())).Should(Succeed())
					Ω(unit.Delete(context.Background())).Should(Succeed())
				})

				It("Should be created", func() {
					Ω(unit.Create(context.Background())).Should(Succeed())
				})
			})
		})
	})

	Describe("Enable", Offset(1), func() {
		var unit god.Unit
		var server *ghttp.Server
		var stdout, stderr *os.File

		BeforeEach(func() {
			server = ghttp.NewServer()

			name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

			// Not easy to have the exact count of requests, so we'll just check if requests were made
			server.SetAllowUnhandledRequests(true)

			logsDir, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			logsDir = path.Join(logsDir, "testdata", "logs")
			Ω(os.MkdirAll(logsDir, 0755)).Should(Succeed())

			f, err := os.CreateTemp(logsDir, fmt.Sprintf("%s-stdout-*.log", name))
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Name()).To(BeARegularFile())

			stdout = f

			f, err = os.CreateTemp(logsDir, fmt.Sprintf("%s-stderr-*.log", name))
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Name()).To(BeARegularFile())

			stderr = f

			requestURI := server.URL() + fmt.Sprintf("/%s", name)

			fmt.Fprintf(GinkgoWriter, "Waiting request on %s...\n", requestURI)

			unit, err = factory(context.Background(), god.Opts().
				WithName(name).
				WithScope(god.ScopeUser).
				WithProgram(CurlPath()).
				WithArguments("-s", requestURI).
				WithUserOwner(os.Getuid()).
				WithInterval(time.Second).
				WithDarwinLimitLoadToSessionType(god.DarwinLimitLoadToSessionBackground).
				WithStandardOutput(stdout.Name()).
				WithErrorOutput(stderr.Name()).
				WithStartLimitInterval(time.Second).
				WithRunAtLoad(true),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(unit).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		AfterEach(func() {
			fmt.Fprintln(GinkgoWriter, "-- UNIT STDOUT --")
			_, err := io.Copy(GinkgoWriter, stdout)
			fmt.Fprintln(GinkgoWriter, "-- END STDOUT ---")
			Expect(err).ToNot(HaveOccurred())

			Ω(stdout.Close()).Should(Succeed())
			Ω(os.Remove(stdout.Name())).Should(Succeed())
		})

		AfterEach(func() {
			fmt.Fprintln(GinkgoWriter, "-- UNIT STDERR --")
			_, err := io.Copy(GinkgoWriter, stderr)
			fmt.Fprintln(GinkgoWriter, "-- END STDERR ---")
			Expect(err).ToNot(HaveOccurred())

			Ω(stderr.Close()).Should(Succeed())
			Ω(os.Remove(stderr.Name())).Should(Succeed())
		})

		Context("A non existing unit", func() {
			PIt("Should return an error", func() {
				Ω(unit.Enable(context.Background())).ShouldNot(Succeed())
			})
		})

		Context("A previously created unit", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			It("Should work", func() {
				Ω(unit.Enable(context.Background())).Should(Succeed())
				// Default throttle interval is 10 seconds
				// Ensure the job is started multiple times.
				Eventually(server.ReceivedRequests).
					WithTimeout(30 * time.Second).
					Should(WithTransform(
						func(requests []*http.Request) int { return len(requests) },
						BeNumerically(">", 1),
					))

				fmt.Fprintf(GinkgoWriter, "%d runs\n", len(server.ReceivedRequests()))
			})

			It("Can be enabled multiple times", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Enable(context.Background())).Should(Succeed())
				Ω(unit.Enable(context.Background())).Should(Succeed())
			})
		})

		Context("A previously deleted unit", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			It("Should return an error", func() {
				Ω(unit.Enable(context.Background())).Should(Succeed())
			})
		})
	})

	Describe("Disable", Offset(1), func() {
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

		Context("A non existing unit", func() {
			PIt("Should return an error", func() {
				Ω(unit.Disable(context.Background())).ShouldNot(Succeed())
			})
		})

		Context("A previously created unit", func() {
			BeforeEach(func() {
				Ω(unit.Enable(context.Background())).Should(Succeed())
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			It("Should work", func() {
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})

			It("Can be enabled multiple times", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Disable(context.Background())).Should(Succeed())
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})
		})

		Context("A previously enabled unit", func() {
			BeforeEach(func() {
				Ω(unit.Enable(context.Background())).Should(Succeed())
				Ω(unit.Create(context.Background())).Should(Succeed())
				Ω(unit.Enable(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			It("Should work", func() {
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})

			It("Can be enabled multiple times", func() {
				currentUser, err := user.Current()
				Expect(err).ToNot(HaveOccurred())

				if currentUser.Uid != "0" {
					Skip("This test requires root privileges")
				}

				Ω(unit.Disable(context.Background())).Should(Succeed())
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})
		})

		Context("A previously deleted unit", func() {
			BeforeEach(func() {
				Ω(unit.Enable(context.Background())).Should(Succeed())
				Ω(unit.Create(context.Background())).Should(Succeed())
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			It("Should return an error", func() {
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})
		})
	})

	Describe("Status", Offset(1), func() {
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

		Context("Of a non existing unit", func() {
			It("Should work", func() {
				status, err := unit.Status(context.Background())
				Expect(err).ToNot(HaveOccurred())

				Ω(status.Exists(context.Background())).Should(BeFalse())
				Ω(status.IsEnabled(context.Background())).Should(BeFalse())
			})
		})

		Context("Of a previously created unit", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			PIt("Should work", func() {
				status, err := unit.Status(context.Background())
				Expect(err).ToNot(HaveOccurred())

				Ω(status.Exists(context.Background())).Should(BeTrue())
				Ω(status.IsEnabled(context.Background())).Should(BeFalse())
			})
		})

		Context("Of an enabled unit", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			PIt("Should work", func() {
				status, err := unit.Status(context.Background())
				Expect(err).ToNot(HaveOccurred())

				Ω(status.Exists(context.Background())).Should(BeTrue())
				Ω(status.IsEnabled(context.Background())).Should(BeFalse())
			})
		})

		Context("Of a disabled unit", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
				Ω(unit.Enable(context.Background())).Should(Succeed())
				Ω(unit.Disable(context.Background())).Should(Succeed())
			})

			AfterEach(func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})

			PIt("Should work", func() {
				status, err := unit.Status(context.Background())
				Expect(err).ToNot(HaveOccurred())

				Ω(status.Exists(context.Background())).Should(BeTrue())
				Ω(status.IsEnabled(context.Background())).Should(BeFalse())
			})
		})
	})

	Describe("Delete", Offset(1), func() {
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

				Ω(unit.Delete(context.Background())).Should(Succeed())
			})
		})

		Context("After installation", func() {
			BeforeEach(func() {
				Ω(unit.Create(context.Background())).Should(Succeed())
			})

			It("should work", func() {
				Ω(unit.Delete(context.Background())).Should(Succeed())
			})
		})
	})
}
