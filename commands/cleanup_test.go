package commands_test

import (
	"github.com/cloudfoundry/bosh-bootloader/commands"
	"github.com/cloudfoundry/bosh-bootloader/fakes"
	"github.com/cloudfoundry/bosh-bootloader/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup", func() {
	Describe("Execute", func() {
		var (
			filter    string
			leftovers *fakes.Leftovers

			cleanup commands.Cleanup
		)

		BeforeEach(func() {
			filter = "banana"
			leftovers = &fakes.Leftovers{}

			cleanup = commands.NewCleanup(leftovers)
		})

		It("calls delete on leftovers with the filter", func() {
			err := cleanup.Execute([]string{"--filter", filter}, storage.State{})
			Expect(err).NotTo(HaveOccurred())

			Expect(leftovers.DeleteCall.CallCount).To(Equal(1))
			Expect(leftovers.DeleteCall.Receives.Filter).To(Equal(filter))
		})

		Context("when parsing flags throws an error", func() {
			It("returns a helpful message", func() {
				err := cleanup.Execute([]string{"--filter"}, storage.State{})
				Expect(err).To(MatchError(ContainSubstring("Parsing clean-up args: flag needs an argument")))

				Expect(leftovers.DeleteCall.CallCount).To(Equal(0))
			})
		})
	})
})
