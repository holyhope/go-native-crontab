package god_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/holyhope/god"
)

var _ = Describe("Options", func() {
	var opts god.Options

	BeforeEach(func() {
		opts = god.Opts()
		Expect(opts).ToNot(BeNil())
	})

	Describe("With no option", func() {
		It("Should not have any options", func() {
			Expect(opts.HasArguments()).To(BeFalse())
			Expect(opts.HasDescription()).To(BeFalse())
			Expect(opts.HasEnvironmentVariables()).To(BeFalse())
			Expect(opts.HasGroupOwner()).To(BeFalse())
			Expect(opts.HasInterval()).To(BeFalse())
			Expect(opts.HasName()).To(BeFalse())
			Expect(opts.HasProgram()).To(BeFalse())
			Expect(opts.HasRunAtLoad()).To(BeFalse())
			Expect(opts.HasScope()).To(BeFalse())
			Expect(opts.HasState()).To(BeFalse())
			Expect(opts.HasUserOwner()).To(BeFalse())
		})
	})

	Describe("With a single option", func() {
		BeforeEach(func() {
			opts = opts.WithName("test")
		})

		It("Contains the option", func() {
			Expect(opts.HasName()).To(BeTrue())
			Expect(opts.Name()).To(Equal("test"))
		})

		It("Does not contain other option", func() {
			Expect(opts.HasArguments()).To(BeFalse())
			Expect(opts.HasDescription()).To(BeFalse())
			Expect(opts.HasEnvironmentVariables()).To(BeFalse())
			Expect(opts.HasGroupOwner()).To(BeFalse())
			Expect(opts.HasInterval()).To(BeFalse())
			// Do not check for opts.HasName()
			Expect(opts.HasProgram()).To(BeFalse())
			Expect(opts.HasRunAtLoad()).To(BeFalse())
			Expect(opts.HasScope()).To(BeFalse())
			Expect(opts.HasState()).To(BeFalse())
			Expect(opts.HasUserOwner()).To(BeFalse())
		})
	})

	Describe("With multiple option", func() {
		BeforeEach(func() {
			opts = opts.
				WithName("test").
				WithProgram("/bin/bash").
				WithScope(god.ScopeSystem)
		})

		It("Contains all the options", func() {
			Expect(opts.Name()).To(Equal("test"))
			Expect(opts.Program()).To(Equal("/bin/bash"))
			Expect(opts.Scope()).To(Equal(god.ScopeSystem))
		})
	})

	Describe("Adding option", func() {
		Context("To an non empty option", func() {
			BeforeEach(func() {
				opts = opts.WithName("test")
			})

			It("Should merge options", func() {
				opts := opts.WithScope(god.ScopeSystem)
				Expect(opts).ToNot(BeNil())
				Expect(opts.Name()).To(Equal("test"))
				Expect(opts.Scope()).To(Equal(god.ScopeSystem))
			})

			It("Should override value with same key", func() {
				opts := opts.WithName("test2")
				Expect(opts).ToNot(BeNil())
				Expect(opts.Name()).To(Equal("test2"))
			})
		})
	})
})
