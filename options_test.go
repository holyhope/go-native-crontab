package god_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/holyhope/god"
)

var _ = Describe("Options", func() {
	Describe("With no option", func() {
		var opts god.Options

		BeforeEach(func() {
			opts = god.With()
		})

		It("Be valid", func() {
			Expect(opts).ToNot(BeNil())
		})

		It("Does not contain other option", func() {
			Expect(opts.Keys()).To(HaveLen(0))
		})
	})
	Describe("With a single option", func() {
		var opts god.Options

		BeforeEach(func() {
			opts = god.With(god.Name, "test")
		})

		It("Contains the option", func() {
			Expect(opts.Get(god.Name)).To(Equal("test"))
		})

		It("Does not contain other option", func() {
			Expect(opts.Keys()).To(HaveLen(1))
		})
	})

	Describe("With multiple option", func() {
		var opts god.Options

		BeforeEach(func() {
			opts = god.With(god.Name, "test", god.Program, "/bin/bash", god.Scope, god.ScopeSystem)
		})

		It("Contains all the options", func() {
			Expect(opts.Get(god.Name)).To(Equal("test"))
			Expect(opts.Get(god.Program)).To(Equal("/bin/bash"))
			Expect(opts.Get(god.Scope)).To(Equal(god.ScopeSystem))
			Expect(opts.Keys()).To(HaveLen(3))
		})
	})

	Describe("Adding option", func() {
		var opts god.Options

		Context("To nil Options", func() {
			BeforeEach(func() {
				opts = nil
			})

			PIt("Should return a new options containing the value", func() {
				opts := opts.And(god.Name, "test")
				Expect(opts).ToNot(BeNil())
				Expect(opts.Get(god.Name)).To(Equal("test"))
			})
		})

		Context("To an non empty option", func() {
			BeforeEach(func() {
				opts = god.With(god.Name, "test")
			})

			It("Should merge options", func() {
				opts := opts.And(god.Scope, god.ScopeSystem)
				Expect(opts).ToNot(BeNil())
				Expect(opts.Get(god.Name)).To(Equal("test"))
				Expect(opts.Get(god.Scope)).To(Equal(god.ScopeSystem))
				Expect(opts.Keys()).To(HaveLen(2))
			})

			It("Should override value with same key", func() {
				opts := opts.And(god.Name, "test2")
				Expect(opts).ToNot(BeNil())
				Expect(opts.Get(god.Name)).To(Equal("test2"))
				Expect(opts.Keys()).To(HaveLen(1))
			})
		})
	})
})
