package terraform_test

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry/bosh-bootloader/fakes"
	"github.com/cloudfoundry/bosh-bootloader/terraform"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TerraformOutputProvider", func() {
	var (
		terraformExecutor       *fakes.TerraformExecutor
		terraformOutputProvider terraform.OutputProvider
	)

	BeforeEach(func() {
		terraformExecutor = &fakes.TerraformExecutor{}
		terraformExecutor.OutputCall.Stub = func(output string) (string, error) {
			switch output {
			case "network_name":
				return "some-network-name", nil
			case "subnetwork_name":
				return "some-subnetwork-name", nil
			case "bosh_open_tag_name":
				return "some-bosh-open-tag-name", nil
			case "internal_tag_name":
				return "some-internal-tag-name", nil
			case "external_ip":
				return "some-external-ip", nil
			case "director_address":
				return "some-director-address", nil
			case "router_backend_service":
				return "some-router-backend-service", nil
			case "ws_lb_ip":
				return "some-ws-lb-ip", nil
			case "ssh_proxy_lb_ip":
				return "some-ssh-proxy-lb-ip", nil
			case "tcp_router_lb_ip":
				return "some-tcp-router-lb-ip", nil
			case "router_lb_ip":
				return "some-router-lb-ip", nil
			case "ssh_proxy_target_pool":
				return "some-ssh-proxy-target-pool", nil
			case "tcp_router_target_pool":
				return "some-tcp-router-target-pool", nil
			case "ws_target_pool":
				return "some-ws-target-pool", nil
			case "concourse_target_pool":
				return "some-concourse-target-pool", nil
			case "concourse_lb_ip":
				return "some-concourse-lb-ip", nil
			case "system_domain_dns_servers":
				return "", errors.New("no dns server exists")
			default:
				return "", nil
			}
		}
		terraformOutputProvider = terraform.NewOutputProvider(terraformExecutor)
	})

	Context("when no lb exists", func() {
		It("returns all terraform outputs except lb related outputs", func() {
			terraformOutputs, err := terraformOutputProvider.Get("some-tf-state", "", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(terraformExecutor.OutputCall.Receives.TFState).To(Equal("some-tf-state"))
			Expect(terraformOutputs).To(Equal(terraform.Outputs{
				ExternalIP:      "some-external-ip",
				NetworkName:     "some-network-name",
				SubnetworkName:  "some-subnetwork-name",
				BOSHTag:         "some-bosh-open-tag-name",
				InternalTag:     "some-internal-tag-name",
				DirectorAddress: "some-director-address",
			}))
		})
	})

	Context("when cf lb exists", func() {
		Context("when the domain is not specified", func() {
			It("returns terraform outputs related to cf lb without system domain DNS servers", func() {
				terraformOutputs, err := terraformOutputProvider.Get("some-tf-state", "cf", false)
				Expect(err).NotTo(HaveOccurred())
				Expect(terraformExecutor.OutputCall.Receives.TFState).To(Equal("some-tf-state"))
				Expect(terraformOutputs).To(Equal(terraform.Outputs{
					ExternalIP:           "some-external-ip",
					NetworkName:          "some-network-name",
					SubnetworkName:       "some-subnetwork-name",
					BOSHTag:              "some-bosh-open-tag-name",
					InternalTag:          "some-internal-tag-name",
					DirectorAddress:      "some-director-address",
					RouterBackendService: "some-router-backend-service",
					SSHProxyTargetPool:   "some-ssh-proxy-target-pool",
					TCPRouterTargetPool:  "some-tcp-router-target-pool",
					WSTargetPool:         "some-ws-target-pool",
					RouterLBIP:           "some-router-lb-ip",
					SSHProxyLBIP:         "some-ssh-proxy-lb-ip",
					TCPRouterLBIP:        "some-tcp-router-lb-ip",
					WebSocketLBIP:        "some-ws-lb-ip",
				}))
			})
		})

		Context("when the domain is specified", func() {
			It("returns terraform outputs related to cf lb with the system domain DNS servers", func() {
				terraformExecutor.OutputCall.Stub = func(output string) (string, error) {
					switch output {
					case "network_name":
						return "some-network-name", nil
					case "subnetwork_name":
						return "some-subnetwork-name", nil
					case "bosh_open_tag_name":
						return "some-bosh-open-tag-name", nil
					case "internal_tag_name":
						return "some-internal-tag-name", nil
					case "external_ip":
						return "some-external-ip", nil
					case "director_address":
						return "some-director-address", nil
					case "router_backend_service":
						return "some-router-backend-service", nil
					case "ws_lb_ip":
						return "some-ws-lb-ip", nil
					case "ssh_proxy_lb_ip":
						return "some-ssh-proxy-lb-ip", nil
					case "tcp_router_lb_ip":
						return "some-tcp-router-lb-ip", nil
					case "router_lb_ip":
						return "some-router-lb-ip", nil
					case "ssh_proxy_target_pool":
						return "some-ssh-proxy-target-pool", nil
					case "tcp_router_target_pool":
						return "some-tcp-router-target-pool", nil
					case "ws_target_pool":
						return "some-ws-target-pool", nil
					case "concourse_target_pool":
						return "some-concourse-target-pool", nil
					case "concourse_lb_ip":
						return "some-concourse-lb-ip", nil
					case "system_domain_dns_servers":
						return "some-name-server-1,\nsome-name-server-2", nil
					default:
						return "", nil
					}
				}

				terraformOutputs, err := terraformOutputProvider.Get("some-tf-state", "cf", true)
				Expect(err).NotTo(HaveOccurred())
				Expect(terraformExecutor.OutputCall.Receives.TFState).To(Equal("some-tf-state"))
				Expect(terraformOutputs).To(Equal(terraform.Outputs{
					ExternalIP:             "some-external-ip",
					NetworkName:            "some-network-name",
					SubnetworkName:         "some-subnetwork-name",
					BOSHTag:                "some-bosh-open-tag-name",
					InternalTag:            "some-internal-tag-name",
					DirectorAddress:        "some-director-address",
					RouterBackendService:   "some-router-backend-service",
					SSHProxyTargetPool:     "some-ssh-proxy-target-pool",
					TCPRouterTargetPool:    "some-tcp-router-target-pool",
					WSTargetPool:           "some-ws-target-pool",
					RouterLBIP:             "some-router-lb-ip",
					SSHProxyLBIP:           "some-ssh-proxy-lb-ip",
					TCPRouterLBIP:          "some-tcp-router-lb-ip",
					WebSocketLBIP:          "some-ws-lb-ip",
					SystemDomainDNSServers: []string{"some-name-server-1", "some-name-server-2"},
				}))
			})
		})
	})

	Context("when concourse lb exists", func() {
		It("returns terraform outputs related to concourse lb", func() {
			terraformOutputs, err := terraformOutputProvider.Get("some-tf-state", "concourse", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(terraformExecutor.OutputCall.Receives.TFState).To(Equal("some-tf-state"))
			Expect(terraformOutputs).To(Equal(terraform.Outputs{
				ExternalIP:          "some-external-ip",
				NetworkName:         "some-network-name",
				SubnetworkName:      "some-subnetwork-name",
				BOSHTag:             "some-bosh-open-tag-name",
				InternalTag:         "some-internal-tag-name",
				DirectorAddress:     "some-director-address",
				ConcourseTargetPool: "some-concourse-target-pool",
				ConcourseLBIP:       "some-concourse-lb-ip",
			}))
		})
	})

	Context("when tfState is empty", func() {
		It("returns an empty terraform outputs", func() {
			terraformOutputs, err := terraformOutputProvider.Get("", "concourse", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(terraformExecutor.OutputCall.CallCount).To(Equal(0))
			Expect(terraformOutputs).To(Equal(terraform.Outputs{}))
		})
	})

	Context("failure cases", func() {
		DescribeTable("returns an error when the outputter fails", func(outputName, lbType string) {
			expectedError := fmt.Sprintf("failed to get %s", outputName)
			terraformExecutor.OutputCall.Stub = func(output string) (string, error) {
				if output == outputName {
					return "", errors.New(expectedError)
				}
				return "", nil
			}

			_, err := terraformOutputProvider.Get(outputName, lbType, true)
			Expect(err).To(MatchError(expectedError))
		},
			Entry("failed to get external_ip", "external_ip", ""),
			Entry("failed to get network_name", "network_name", ""),
			Entry("failed to get subnetwork_name", "subnetwork_name", ""),
			Entry("failed to get bosh_open_tag_name", "bosh_open_tag_name", ""),
			Entry("failed to get internal_tag_name", "internal_tag_name", ""),
			Entry("failed to get director_address", "director_address", ""),
			Entry("failed to get router_backend_service", "router_backend_service", "cf"),
			Entry("failed to get ssh_proxy_target_pool", "ssh_proxy_target_pool", "cf"),
			Entry("failed to get tcp_router_target_pool", "tcp_router_target_pool", "cf"),
			Entry("failed to get ws_target_pool", "ws_target_pool", "cf"),
			Entry("failed to get router_lb_ip", "router_lb_ip", "cf"),
			Entry("failed to get ssh_proxy_lb_ip", "ssh_proxy_lb_ip", "cf"),
			Entry("failed to get tcp_router_lb_ip", "tcp_router_lb_ip", "cf"),
			Entry("failed to get ws_lb_ip", "ws_lb_ip", "cf"),
			Entry("failed to get concourse_target_pool", "concourse_target_pool", "concourse"),
			Entry("failed to get concourse_lb_ip", "concourse_lb_ip", "concourse"),
			Entry("failed to get system_domain_dns_servers", "system_domain_dns_servers", "cf"),
		)
	})
})