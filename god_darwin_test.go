package god_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "embed"

	crontab "github.com/holyhope/god"
	god "github.com/holyhope/god"
	"github.com/iancoleman/strcase"
)

var _ = Describe("CrontabDarwin", func() {
	Describe("New unit", func() {
		var opts []god.FactoryOpts

		BeforeEach(func() {
			opts = []god.FactoryOpts{}
		})

		Context("With no options", func() {
			It("should return a unit", func() {
				var err error
				unit, err := god.New(context.Background(), opts...)
				Expect(err).ToNot(HaveOccurred())
				Expect(unit).ToNot(BeNil())
			})
		})
	})

	Describe("Install", func() {
		var unit god.Unit
		var opts []god.FactoryOpts

		BeforeEach(func() {
			opts = []god.FactoryOpts{}
		})

		JustBeforeEach(func() {
			var err error
			unit, err = god.New(context.Background(), opts...)
			Expect(err).ToNot(HaveOccurred())
			Expect(unit).ToNot(BeNil())
		})

		Context("with no name", func() {
			It("should fail", func() {
				Ω(unit.Install(context.Background())).Should(MatchError(&god.MissingOptionsError{
					Missings: []string{"UnitName", "Scope"},
				}))
			})
		})

		Context("With name", func() {
			BeforeEach(func() {
				fileName := fmt.Sprintf("com.github.holyhope.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))
				opts = append(opts, god.UnitName(fileName))
			})

			Context("with no scope", func() {
				It("should fail", func() {
					Ω(unit.Install(context.Background())).Should(MatchError(&god.MissingOptionsError{
						Missings: []string{"Scope"},
					}))
				})
			})

			Context("with scope", func() {
				BeforeEach(func() {
					opts = append(opts, crontab.ScopeUser)
				})

				Context("Without entries", func() {
					It("should work", func() {
						Ω(unit.Install(context.Background())).Should(Succeed())
					})
				})
			})
		})
	})
})
