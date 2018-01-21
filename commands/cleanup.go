package commands

import (
	"fmt"

	"github.com/cloudfoundry/bosh-bootloader/flags"
	"github.com/cloudfoundry/bosh-bootloader/storage"
)

type Cleanup struct {
	leftovers Leftovers
}

type Leftovers interface {
	Delete(string) error
}

func NewCleanup(leftovers Leftovers) Cleanup {
	return Cleanup{
		leftovers: leftovers,
	}
}

func (c Cleanup) CheckFastFails(subcommandFlags []string, state storage.State) error {
	return nil
}

func (c Cleanup) Execute(args []string, state storage.State) error {
	var filter string
	f := flags.New("clean-up")
	f.String(&filter, "filter", "")

	err := f.Parse(args)
	if err != nil {
		return fmt.Errorf("Parsing clean-up args: %s", err)
	}

	return c.leftovers.Delete(filter)
}
