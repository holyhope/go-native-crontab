package launchd_test

import (
	"context"
	"fmt"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	_ "embed"

	"github.com/holyhope/god"
	"github.com/holyhope/god/launchd"

	"github.com/iancoleman/strcase"
)

//go:embed testdata/test.plist
var plistContent []byte

var _ = Describe("Launchd", func() {
	Describe("New unit", func() {
		var opts god.Options

		BeforeEach(func() {
			opts = nil
		})

		Context("With no options", func() {
			It("should return an error", func() {
				_, err := launchd.New(context.Background(), opts)
				Expect(err).To(Or(
					MatchError(&god.InvalidOptionError{
						Key:   god.Name,
						Value: nil,
					}),
					MatchError(&god.InvalidOptionError{
						Key:   god.Scope,
						Value: nil,
					}),
				))
			})
		})
	})

	Describe("Install", func() {
		var unit god.Unit

		Context("With name", func() {
			BeforeEach(func() {
				name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

				var err error
				unit, err = launchd.New(context.Background(), god.With(
					god.Name, name,
					god.Scope, god.ScopeUser,
					god.Program, "bash",
					god.ProgramArguments, []string{"-c", `echo 'Hello, world!'`},
				))
				Expect(err).ToNot(HaveOccurred())
				Expect(unit).ToNot(BeNil())
			})

			It("should work", func() {
				Ω(unit.Install(context.Background())).Should(Succeed())
				Ω(unit.Uninstall(context.Background())).Should(Succeed())
			})
		})
	})

	Describe("Uninstall", func() {
		var unit god.Unit

		BeforeEach(func() {
			name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

			var err error
			unit, err = launchd.New(context.Background(), god.With(
				god.Name, name,
				god.Scope, god.ScopeUser,
				god.Program, "bash",
				god.ProgramArguments, []string{"-c", `echo 'Hello, world!'`},
			))
			Expect(err).ToNot(HaveOccurred())
			Expect(unit).ToNot(BeNil())
		})

		Context("without installation", func() {
			It("should fail", func() {
				Ω(unit.Uninstall(context.Background())).ShouldNot(Succeed())
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

	Describe("ToPlist", func() {
		type PList interface {
			ToPlist(io.Writer) error
		}

		var unit PList

		BeforeEach(func() {
			name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

			u, err := launchd.New(
				context.Background(),
				god.With(
					god.Name, name,
					god.Scope, god.ScopeUser,
					god.Program, "bash",
					god.ProgramArguments, []string{"-c", `echo 'Hello, world!'`},
				),
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(u).ShouldNot(BeNil())

			unit = u.(PList)
		})

		It("should return proper plist content", func() {
			reader, writer := io.Pipe()
			defer reader.Close()

			go func(writer io.WriteCloser) {
				defer GinkgoRecover()
				defer writer.Close()

				Ω(unit.ToPlist(writer)).Should(Succeed())
			}(writer)

			buffer := gbytes.BufferReader(reader)
			Eventually(buffer).Should(WithTransform(func(buffer *gbytes.Buffer) []byte {
				return buffer.Contents()
			}, MatchXML(plistContent)))
		})
	})
})
