package launchd_test

import (
	"context"
	"fmt"
	"io"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/holyhope/god"
	"github.com/holyhope/god/internal/tests"
	"github.com/holyhope/god/launchd"
	"github.com/iancoleman/strcase"
)

//go:embed testdata/test.plist
var plistContent []byte

var _ = Describe("Launchd", func() {
	tests.NewSuite(launchd.New)

	Describe("ToPlist", func() {
		type PList interface {
			ToPlist(io.Writer) error
		}

		var unit PList

		BeforeEach(func() {
			name := fmt.Sprintf("com.github.holyhope.god.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))

			u, err := launchd.New(context.Background(), god.Opts().
				WithName(name).
				WithScope(god.ScopeUser).
				WithProgram("/bin/bash").
				WithArguments("-c", `echo 'Hello, world!'`).
				WithUserOwner(0).
				WithDarwinLimitLoadToSessionType(god.DarwinLimitLoadToSessionBackground),
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
