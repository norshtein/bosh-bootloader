package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	client       instancesClient
	resourceTags resourceTags
	id           *string
	identifier   string
	rtype        string
}

func NewInstance(client instancesClient, resourceTags resourceTags, id, keyName *string, tags []*awsec2.Tag) Instance {
	identifier := *id

	extra := []string{}
	for _, t := range tags {
		extra = append(extra, fmt.Sprintf("%s:%s", *t.Key, *t.Value))
	}

	if keyName != nil && *keyName != "" {
		extra = append(extra, fmt.Sprintf("KeyPairName:%s", *keyName))
	}

	if len(extra) > 0 {
		identifier = fmt.Sprintf("%s (%s)", *id, strings.Join(extra, ", "))
	}

	return Instance{
		client:       client,
		resourceTags: resourceTags,
		id:           id,
		identifier:   identifier,
		rtype:        "EC2 Instance",
	}
}

func (i Instance) Delete() error {
	addresses, err := i.client.DescribeAddresses(&awsec2.DescribeAddressesInput{
		Filters: []*awsec2.Filter{{
			Name:   aws.String("instance-id"),
			Values: []*string{i.id},
		}},
	})
	if err != nil {
		return fmt.Errorf("Describe addresses: %s", err)
	}

	_, err = i.client.TerminateInstances(&awsec2.TerminateInstancesInput{InstanceIds: []*string{i.id}})
	if err != nil {
		return fmt.Errorf("Terminate: %s", err)
	}

	err = i.resourceTags.Delete("instance", *i.id)
	if err != nil {
		return fmt.Errorf("Delete resource tags: %s", err)
	}

	for _, a := range addresses.Addresses {
		_, err = i.client.ReleaseAddress(&awsec2.ReleaseAddressInput{AllocationId: a.AllocationId})
		if err != nil {
			return fmt.Errorf("Release address: %s", err)
		}
	}

	return nil
}

func (i Instance) Name() string {
	return i.identifier
}

func (i Instance) Type() string {
	return i.rtype
}
