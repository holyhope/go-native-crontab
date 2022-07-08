package crontab_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "embed"

	crontab "github.com/holyhope/go-native-crontab"
	"github.com/iancoleman/strcase"
)

var _ = Describe("CrontabDarwin", func() {
	var cronTab crontab.CronTab

	BeforeEach(func() {
		var err error

		cronTab, err = crontab.New(context.Background())
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Add entry", func() {
		It("should work", func() {
			Ω(cronTab.Add(context.Background(), time.Second, "test", "test")).Should(Succeed())
		})
	})

	Describe("Install", func() {
		var opts []crontab.InstallOpts

		BeforeEach(func() {
			opts = []crontab.InstallOpts{}
		})

		Context("with no filename", func() {
			It("should fail", func() {
				_, err := cronTab.Install(context.Background(), opts...)
				Ω(err).Should(MatchError(&crontab.MissingOptionsError{
					Name: "FileName",
				}))
			})
		})

		Context("With filename", func() {
			BeforeEach(func() {
				fileName := fmt.Sprintf("com.github.holyhope.test.%s", strcase.ToSnake(CurrentSpecReport().FullText()))
				opts = append(opts, crontab.FileName(fileName))
			})

			Context("with no scope", func() {
				It("should fail", func() {
					_, err := cronTab.Install(context.Background(), opts...)
					Ω(err).Should(MatchError(&crontab.MissingOptionsError{
						Name: "Scope",
					}))
				})
			})

			Context("with scope", func() {
				BeforeEach(func() {
					opts = append(opts, crontab.UserScope)
				})

				Context("Without entries", func() {
					It("should work", func() {
						ict, err := cronTab.Install(context.Background(), opts...)
						Expect(err).ToNot(HaveOccurred())
						Ω(ict.Uninstall(context.Background())).Should(Succeed())
					})
				})

				Context("With a single entry", func() {
					BeforeEach(func() {
						Ω(cronTab.Add(context.Background(), time.Hour, "test", "test")).Should(Succeed())
					})

					It("should work", func() {
						cronTab, err := cronTab.Install(context.Background(), opts...)
						Expect(err).ToNot(HaveOccurred())
						Ω(cronTab.Uninstall(context.Background())).Should(Succeed())
					})
				})

				Context("With multiple entries", func() {
					BeforeEach(func() {
						Ω(cronTab.Add(context.Background(), time.Second, "test1", "arg1")).Should(Succeed())
						Ω(cronTab.Add(context.Background(), time.Minute, "test2", "arg2")).Should(Succeed())
						Ω(cronTab.Add(context.Background(), time.Minute, "test3", "arg3")).Should(Succeed())
					})

					It("should work", func() {
						ict, err := cronTab.Install(context.Background(), opts...)
						Expect(err).ToNot(HaveOccurred())
						Ω(ict.Uninstall(context.Background())).Should(Succeed())
					})
				})
			})
		})
	})
})
