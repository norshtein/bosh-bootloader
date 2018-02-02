package commands

import (
	"fmt"

	"github.com/cloudfoundry/bosh-bootloader/flags"
	"github.com/cloudfoundry/bosh-bootloader/storage"

	awsleftovers "github.com/genevievelesperance/leftovers/aws"
	azureleftovers "github.com/genevievelesperance/leftovers/azure"
	gcpleftovers "github.com/genevievelesperance/leftovers/gcp"
)

type CleanupLeftovers struct {
	logger      logger
	iaas        string
	credentials []string
}

func NewCleanupLeftovers(logger logger, iaas string, credentials ...string) CleanupLeftovers {
	return CleanupLeftovers{
		logger:      logger,
		iaas:        iaas,
		credentials: credentials,
	}
}

func (l CleanupLeftovers) CheckFastFails(subcommandFlags []string, state storage.State) error {
	return nil
}

func (l CleanupLeftovers) Execute(subcommandFlags []string, state storage.State) error {
	var (
		deleter interface {
			Delete(string) error
		}
		err error
	)

	switch l.iaas {
	case "aws":
		deleter, err = awsleftovers.NewLeftovers(l.logger, l.credentials...)
	case "azure":
		deleter, err = azureleftovers.NewLeftovers(l.logger, l.credentials...)
	case "gcp":
		deleter, err = gcpleftovers.NewLeftovers(l.logger, l.credentials...)
	}
	if err != nil {
		return fmt.Errorf("Collecting leftovers: %s", err)
	}

	var filter string
	f := flags.New("cleanup-leftovers")
	f.String(&filter, "filter", "")

	err = f.Parse(subcommandFlags)
	if err != nil {
		return fmt.Errorf("Parsing cleanup-leftovers args: %s", err)
	}

	return deleter.Delete(filter)
}

func (l CleanupLeftovers) Usage() string {
	return fmt.Sprintf("%s%s%s", CleanupLeftoversCommandUsage, requiresCredentials, Credentials)
}
