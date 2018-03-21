package storage_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/bosh-bootloader/fakes"
	"github.com/cloudfoundry/bosh-bootloader/storage"
	uuid "github.com/nu7hatch/gouuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var (
		fileIO  *fakes.FileIO
		store   storage.Store
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")

		fileIO = &fakes.FileIO{}

		store = storage.NewStore(tempDir, fileIO)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		storage.ResetUUIDNewV4()
		storage.ResetMarshalIndent()
	})

	Describe("Set", func() {
		Context("when credhub is enabled", func() {
			It("stores the state into a file, without IAAS credentials", func() {
				storage.SetUUIDNewV4(func() (*uuid.UUID, error) {
					return &uuid.UUID{
						0x01, 0x02, 0x03, 0x04,
						0x05, 0x06, 0x07, 0x08,
						0x09, 0x10, 0x11, 0x12,
						0x13, 0x14, 0x15, 0x16}, nil
				})
				err := store.Set(storage.State{
					BBLVersion: "5.3.0",
					IAAS:       "aws",
					AWS: storage.AWS{
						AccessKeyID:     "some-aws-access-key-id",
						SecretAccessKey: "some-aws-secret-access-key",
						Region:          "some-region",
					},
					Azure: storage.Azure{
						ClientID:       "client-id",
						ClientSecret:   "client-secret",
						Region:         "some-azure-region",
						SubscriptionID: "subscription-id",
						TenantID:       "tenant-id",
					},
					GCP: storage.GCP{
						ServiceAccountKey: "some-service-account-key",
						ProjectID:         "some-project-id",
						Zone:              "some-zone",
						Region:            "some-region",
						Zones:             []string{"some-zone", "some-other-zone"},
					},
					VSphere: storage.VSphere{
						VCenterUser:     "user",
						VCenterPassword: "password",
						VCenterIP:       "ip",
						VCenterDC:       "dc",
						VCenterCluster:  "cluster",
						VCenterRP:       "rp",
						Network:         "network",
						VCenterDS:       "ds",
						Subnet:          "10.0.0.0/24",
					},
					OpenStack: storage.OpenStack{
						InternalCidr:         "cidr",
						ExternalIP:           "ext-ip",
						AuthURL:              "auth-url",
						AZ:                   "az",
						DefaultKeyName:       "default-key-name",
						DefaultSecurityGroup: "default-security-group",
						NetworkID:            "network-id",
						Password:             "password",
						Username:             "username",
						Project:              "project",
						Domain:               "domain",
						Region:               "region",
					},
					LB: storage.LB{
						Type:   "some-type",
						Cert:   "some-cert",
						Key:    "some-key",
						Chain:  "some-chain",
						Domain: "some-domain",
					},
					Jumpbox: storage.Jumpbox{
						URL:       "some-jumpbox-url",
						Manifest:  "name: jumpbox",
						Variables: "some-jumpbox-vars",
						State: map[string]interface{}{
							"key": "value",
						},
					},
					BOSH: storage.BOSH{
						DirectorName:           "some-director-name",
						DirectorUsername:       "some-director-username",
						DirectorPassword:       "some-director-password",
						DirectorAddress:        "some-director-address",
						DirectorSSLCA:          "some-bosh-ssl-ca",
						DirectorSSLCertificate: "some-bosh-ssl-certificate",
						DirectorSSLPrivateKey:  "some-bosh-ssl-private-key",
						State: map[string]interface{}{
							"key": "value",
						},
						Variables: "some-vars",
						Manifest:  "name: bosh",
					},
					EnvID:   "some-env-id",
					TFState: "some-tf-state",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(fileIO.WriteFileCall.Receives[0].Filename).To(Equal(filepath.Join(tempDir, "bbl-state.json")))
				Expect(fileIO.WriteFileCall.Receives[0].Mode).To(Equal(os.FileMode(0644)))
				Expect(fileIO.WriteFileCall.Receives[0].Contents).To(MatchJSON(`{
				"version": 14,
				"bblVersion": "5.3.0",
				"iaas": "aws",
				"id": "01020304-0506-0708-0910-111213141516",
				"envID": "some-env-id",
				"noDirector": false,
				"aws": {
					"region": "some-region"
				},
				"azure": {
					"region": "some-azure-region"
				},
				"gcp": {
					"zone": "some-zone",
					"region": "some-region",
					"zones": ["some-zone", "some-other-zone"]
				},
				"vsphere": {},
				"openstack": {},
				"lb": {
					"type": "some-type",
					"cert": "some-cert",
					"key": "some-key",
					"chain": "some-chain",
					"domain": "some-domain"
				},
				"jumpbox":{
					"url": "some-jumpbox-url",
					"variables": "some-jumpbox-vars",
					"manifest": "name: jumpbox",
					"state": {
						"key": "value"
					}
				},
				"bosh":{
					"directorName": "some-director-name",
					"directorUsername": "some-director-username",
					"directorPassword": "some-director-password",
					"directorAddress": "some-director-address",
					"directorSSLCA": "some-bosh-ssl-ca",
					"directorSSLCertificate": "some-bosh-ssl-certificate",
					"directorSSLPrivateKey": "some-bosh-ssl-private-key",
					"variables":   "some-vars",
					"manifest": "name: bosh",
					"state": {
						"key": "value"
					}
				},
				"tfState": "some-tf-state",
				"latestTFOutput": ""
		    	}`))
			})
		})

		Context("when the state is empty", func() {
			It("removes the bbl-state.json file", func() {
				err := store.Set(storage.State{})
				Expect(err).NotTo(HaveOccurred())

				Expect(fileIO.RemoveCall.Receives[0].Name).To(Equal(filepath.Join(tempDir, "bbl-state.json")))
			})

			It("removes bosh *-env scripts", func() {
				createDirector := filepath.Join(tempDir, "create-director.sh")
				createJumpbox := filepath.Join(tempDir, "create-jumpbox.sh")
				deleteDirector := filepath.Join(tempDir, "delete-director.sh")
				deleteJumpbox := filepath.Join(tempDir, "delete-jumpbox.sh")

				err := store.Set(storage.State{})
				Expect(err).NotTo(HaveOccurred())

				Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: createDirector}))
				Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: deleteDirector}))
				Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: deleteJumpbox}))
				Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: createJumpbox}))
			})

			DescribeTable("removing bbl-created directories",
				func(directory string, expectToBeDeleted bool) {
					err := store.Set(storage.State{})
					Expect(err).NotTo(HaveOccurred())

					if expectToBeDeleted {
						Expect(fileIO.RemoveAllCall.Receives).To(ContainElement(fakes.RemoveAllReceive{
							Path: filepath.Join(tempDir, directory),
						}))
					} else {
						Expect(fileIO.RemoveAllCall.Receives).NotTo(ContainElement(fakes.RemoveAllReceive{
							Path: filepath.Join(tempDir, directory),
						}))
					}
				},
				Entry(".terraform", ".terraform", true),
				Entry("bosh-deployment", "bosh-deployment", true),
				Entry("jumpbox-deployment", "jumpbox-deployment", true),
				Entry("bbl-ops-files", "bbl-ops-files", true),
				Entry("non-bbl directory", "foo", false),
			)

			Describe("cloud-config", func() {
				var (
					cloudConfigBase string
					cloudConfigOps  string
				)
				BeforeEach(func() {
					cloudConfigBase = filepath.Join(tempDir, "cloud-config", "cloud-config.yml")
					cloudConfigOps = filepath.Join(tempDir, "cloud-config", "ops.yml")
				})

				It("removes the ops file, base file, and directory", func() {
					err := store.Set(storage.State{})
					Expect(err).NotTo(HaveOccurred())

					Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: cloudConfigBase}))
					Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: cloudConfigOps}))
					Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{
						Name: filepath.Join(tempDir, "cloud-config"),
					}))
				})
			})

			Describe("vars", func() {
				Context("when the vars directory contains only bbl files", func() {
					BeforeEach(func() {
						fileIO.ReadDirCall.Returns.FileInfos = []os.FileInfo{
							fakes.FileInfo{FileName: "bbl.tfvars"},
							fakes.FileInfo{FileName: "bosh-state.json"},
							fakes.FileInfo{FileName: "director-vars-file.yml"},
							fakes.FileInfo{FileName: "director-vars-store.yml"},
							fakes.FileInfo{FileName: "jumpbox-state.json"},
							fakes.FileInfo{FileName: "jumpbox-vars-file.yml"},
							fakes.FileInfo{FileName: "jumpbox-vars-store.yml"},
							fakes.FileInfo{FileName: "terraform.tfstate"},
							fakes.FileInfo{FileName: "terraform.tfstate.backup"},
						}
					})

					It("removes the directory", func() {
						err := store.Set(storage.State{})
						Expect(err).NotTo(HaveOccurred())

						Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{
							Name: filepath.Join(tempDir, "vars", "bbl.tfvars"),
						}))
						Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{
							Name: filepath.Join(tempDir, "vars"),
						}))
					})
				})

				Context("when the vars directory contains user managed files", func() {
					BeforeEach(func() {
						fileIO.ReadDirCall.Returns.FileInfos = []os.FileInfo{
							fakes.FileInfo{FileName: "user-managed-file"},
							fakes.FileInfo{FileName: "terraform.tfstate.backup"},
						}
					})

					It("spares user managed files", func() {
						err := store.Set(storage.State{})
						Expect(err).NotTo(HaveOccurred())

						Expect(fileIO.RemoveCall.Receives).NotTo(ContainElement(fakes.RemoveReceive{
							Name: filepath.Join(tempDir, "vars", "user-managed-file"),
						}))
					})
				})
			})

			Describe("terraform", func() {
				It("removes the bbl template and directory", func() {
					bblTerraformTemplate := filepath.Join(tempDir, "terraform", "bbl-template.tf")

					err := store.Set(storage.State{})
					Expect(err).NotTo(HaveOccurred())

					Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{Name: bblTerraformTemplate}))
					Expect(fileIO.RemoveCall.Receives).To(ContainElement(fakes.RemoveReceive{
						Name: filepath.Join(tempDir, "terraform"),
					}))
				})
			})

			Context("when the bbl-state.json file does not exist", func() {
				It("does nothing", func() {
					err := store.Set(storage.State{})
					Expect(err).NotTo(HaveOccurred())

					Expect(len(fileIO.WriteFileCall.Receives)).To(Equal(0))
				})
			})

			Context("failure cases", func() {
				Context("when the bbl-state.json file cannot be removed", func() {
					BeforeEach(func() {
						fileIO.RemoveCall.Returns = []fakes.RemoveReturn{{Error: errors.New("permission denied")}}
					})

					It("returns an error", func() {
						err := store.Set(storage.State{})
						Expect(err).To(MatchError(ContainSubstring("permission denied")))
					})
				})

				Context("when uuid new V4 fails", func() {
					It("returns an error", func() {
						storage.SetUUIDNewV4(func() (*uuid.UUID, error) {
							return nil, errors.New("some error")
						})
						err := store.Set(storage.State{
							IAAS: "some-iaas",
						})
						Expect(err).To(MatchError("Create state ID: some error"))
					})
				})
			})
		})

		Context("failure cases", func() {
			Context("when json marshalling fails", func() {
				BeforeEach(func() {
					storage.SetMarshalIndent(func(state interface{}, prefix string, indent string) ([]byte, error) {
						return []byte{}, errors.New("failed to marshal JSON")
					})
				})

				It("returns an error", func() {
					err := store.Set(storage.State{
						IAAS: "aws",
					})
					Expect(err).To(MatchError("failed to marshal JSON"))
				})
			})

			Context("when the directory does not exist", func() {
				BeforeEach(func() {
					fileIO.StatCall.Returns.Error = errors.New("no such file or directory")
				})

				It("returns an error", func() {
					store = storage.NewStore("non-valid-dir", fileIO)
					err := store.Set(storage.State{})
					Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				})
			})

			Context("when it fails to open the bbl-state.json file", func() {
				BeforeEach(func() {
					fileIO.WriteFileCall.Returns = []fakes.WriteFileReturn{{Error: errors.New("permission denied")}}
				})

				It("returns an error", func() {
					err := store.Set(storage.State{EnvID: "something"})
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})

	DescribeTable("get dirs returns the path to an existing directory",
		func(subdirectory string, getDirsFunc func() (string, error)) {
			expectedDir := filepath.Join(tempDir, subdirectory)

			actualDir, err := getDirsFunc()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualDir).To(Equal(expectedDir))

			if len(subdirectory) > 0 {
				Expect(fileIO.MkdirAllCall.Receives.Dir).To(Equal(expectedDir))
			}
		},
		Entry("cloud-config", "cloud-config", func() (string, error) { return store.GetCloudConfigDir() }),
		Entry("state", "", func() (string, error) { return store.GetStateDir(), nil }),
		Entry("vars", "vars", func() (string, error) { return store.GetVarsDir() }),
		Entry("terraform", "terraform", func() (string, error) { return store.GetTerraformDir() }),
		Entry("bosh-deployment", "bosh-deployment", func() (string, error) { return store.GetDirectorDeploymentDir() }),
		Entry("jumpbox-deployment", "jumpbox-deployment", func() (string, error) { return store.GetJumpboxDeploymentDir() }),
	)

	DescribeTable("get dirs returns an error when the subdirectory cannot be created",
		func(subdirectory string, getDirsFunc func() (string, error)) {
			expectedDir := filepath.Join(tempDir, subdirectory)
			fileIO.MkdirAllCall.Returns.Error = errors.New("not a directory")

			_, err := getDirsFunc()
			Expect(err).To(MatchError(ContainSubstring("not a directory")))
			Expect(fileIO.MkdirAllCall.Receives.Dir).To(Equal(expectedDir))
		},
		Entry("cloud-config", "cloud-config", func() (string, error) { return store.GetCloudConfigDir() }),
		Entry("vars", "vars", func() (string, error) { return store.GetVarsDir() }),
		Entry("terraform", "terraform", func() (string, error) { return store.GetTerraformDir() }),
		Entry("bosh-deployment", "bosh-deployment", func() (string, error) { return store.GetDirectorDeploymentDir() }),
		Entry("jumpbox-deployment", "jumpbox-deployment", func() (string, error) { return store.GetJumpboxDeploymentDir() }),
	)

	Describe("GetCloudConfigDir", func() {
		var expectedCloudConfigPath string

		BeforeEach(func() {
			expectedCloudConfigPath = filepath.Join(tempDir, "cloud-config")
		})

		Context("if the cloud-config subdirectory exists", func() {
			It("returns the path to the cloud-config directory", func() {
				cloudConfigDir, err := store.GetCloudConfigDir()
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudConfigDir).To(Equal(expectedCloudConfigPath))
			})
		})

		Context("failure cases", func() {
			Context("when there is a name collision with an existing file", func() {
				BeforeEach(func() {
					fileIO.MkdirAllCall.Returns.Error = errors.New("not a directory")
				})

				It("returns an error", func() {
					cloudConfigDir, err := store.GetCloudConfigDir()
					Expect(err).To(MatchError(ContainSubstring("not a directory")))
					Expect(cloudConfigDir).To(Equal(""))
				})
			})
		})
	})
})
