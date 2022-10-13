package discovery

import (
	"os"
	"testing"
	"time"

	"github.com/barcostreams/barco/internal/conf"
	"github.com/barcostreams/barco/internal/test/conf/mocks"
	dbMocks "github.com/barcostreams/barco/internal/test/localdb/mocks"
	. "github.com/barcostreams/barco/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Discovery Suite")
}

var _ = Describe("discoverer", func() {
	AfterEach(func() {
		os.Setenv(envBrokerNames, "")
		os.Setenv(envOrdinal, "")
	})

	Describe("Init()", func() {
		It("should parse fixed brokers", func() {
			os.Setenv(envOrdinal, "2")
			os.Setenv(envBrokerNames, "abc,def,ghi")
			d := &discoverer{
				localDb: newLocalDbWithNoRecords(),
				config:  newConfigFake(2),
			}

			defer d.Close()
			d.Init()

			Expect(d.Topology().Brokers).To(Equal([]BrokerInfo{
				{
					IsSelf:   false,
					Ordinal:  0,
					HostName: "abc",
				}, {
					IsSelf:   false,
					Ordinal:  1,
					HostName: "def",
				}, {
					IsSelf:   true,
					Ordinal:  2,
					HostName: "ghi",
				},
			}))
		})

		It("should set topology for 3 desired replicas", func() {
			d := &discoverer{
				config: &configFake{
					ordinal:      1,
					baseHostName: "barco-",
				},
				k8sClient: &k8sClientFake{3},
				localDb:   newLocalDbWithNoRecords(),
			}

			defer d.Close()
			d.Init()

			Expect(d.Topology().Brokers).To(Equal([]BrokerInfo{
				{IsSelf: false, Ordinal: 0, HostName: "barco-0.barco"},
				{IsSelf: true, Ordinal: 1, HostName: "barco-1.barco"},
				{IsSelf: false, Ordinal: 2, HostName: "barco-2.barco"},
			}))
			Expect(d.Topology().LocalIndex).To(Equal(BrokerIndex(1)))
		})

		It("should parse 6 real brokers", func() {
			d := &discoverer{
				config: &configFake{
					ordinal:      2,
					baseHostName: "barco-",
				},
				k8sClient: &k8sClientFake{6},
				localDb:   newLocalDbWithNoRecords(),
			}

			defer d.Close()
			d.Init()

			Expect(d.Topology().Brokers).To(Equal([]BrokerInfo{
				{IsSelf: false, Ordinal: 0, HostName: "barco-0.barco"},
				{IsSelf: false, Ordinal: 3, HostName: "barco-3.barco"},
				{IsSelf: false, Ordinal: 1, HostName: "barco-1.barco"},
				{IsSelf: false, Ordinal: 4, HostName: "barco-4.barco"},
				{IsSelf: true, Ordinal: 2, HostName: "barco-2.barco"},
				{IsSelf: false, Ordinal: 5, HostName: "barco-5.barco"},
			}))

			Expect(d.Topology().LocalIndex).To(Equal(BrokerIndex(4)))
		})
	})

	Describe("Topology().Peers()", func() {
		It("should return all brokers except self", func() {
			d := &discoverer{
				config: &configFake{
					ordinal:      1,
					baseHostName: "barco-",
				},
				k8sClient: &k8sClientFake{3},
				localDb:   newLocalDbWithNoRecords(),
			}

			defer d.Close()
			d.Init()

			Expect(d.Topology().Peers()).To(Equal([]BrokerInfo{
				{
					IsSelf:   false,
					Ordinal:  0,
					HostName: "barco-0.barco",
				}, {
					IsSelf:   false,
					Ordinal:  2,
					HostName: "barco-2.barco",
				},
			}))
		})
	})

	Describe("createTopology()", func() {
		It("should return the brokers in placement order for 3 broker cluster", func() {
			config := new(mocks.Config)
			config.On("BaseHostName").Return("barco-")
			config.On("ServiceName").Return("svc")
			config.On("PodNamespace").Return("streams")
			config.On("Ordinal").Return(1)

			topology := createTopology(3, config)
			Expect(topology.Brokers).To(Equal([]BrokerInfo{
				{IsSelf: false, Ordinal: 0, HostName: "barco-0.svc.streams"},
				{IsSelf: true, Ordinal: 1, HostName: "barco-1.svc.streams"},
				{IsSelf: false, Ordinal: 2, HostName: "barco-2.svc.streams"},
			}))
			Expect(topology.LocalIndex).To(Equal(BrokerIndex(1)))
		})

		It("should return the brokers in placement order for 6 broker cluster", func() {
			config := new(mocks.Config)
			config.On("BaseHostName").Return("broker-")
			config.On("ServiceName").Return("svc2")
			config.On("PodNamespace").Return("")
			config.On("Ordinal").Return(1)

			topology := createTopology(6, config)
			Expect(topology.Brokers).To(Equal([]BrokerInfo{
				{IsSelf: false, Ordinal: 0, HostName: "broker-0.svc2"},
				{IsSelf: false, Ordinal: 3, HostName: "broker-3.svc2"},
				{IsSelf: true, Ordinal: 1, HostName: "broker-1.svc2"},
				{IsSelf: false, Ordinal: 4, HostName: "broker-4.svc2"},
				{IsSelf: false, Ordinal: 2, HostName: "broker-2.svc2"},
				{IsSelf: false, Ordinal: 5, HostName: "broker-5.svc2"},
			}))
			Expect(topology.LocalIndex).To(Equal(BrokerIndex(2)))
			Expect(topology.GetIndex(0)).To(Equal(BrokerIndex(0)))
			Expect(topology.GetIndex(3)).To(Equal(BrokerIndex(1)))
			Expect(topology.GetIndex(1)).To(Equal(BrokerIndex(2)))
			Expect(topology.GetIndex(4)).To(Equal(BrokerIndex(3)))
			Expect(topology.GetIndex(2)).To(Equal(BrokerIndex(4)))
			Expect(topology.GetIndex(5)).To(Equal(BrokerIndex(5)))

			config = new(mocks.Config)
			config.On("BaseHostName").Return("broker-")
			config.On("ServiceName").Return("svc2")
			config.On("PodNamespace").Return("")
			config.On("Ordinal").Return(2)
			topology = createTopology(6, config)
			Expect(topology.LocalIndex).To(Equal(BrokerIndex(4)))
		})

		It("should return the brokers in placement order for 12 broker cluster", func() {
			config := new(mocks.Config)
			config.On("BaseHostName").Return("broker-")
			config.On("ServiceName").Return("barco")
			config.On("PodNamespace").Return("")
			config.On("Ordinal").Return(4)

			topology := createTopology(12, config)
			Expect(topology.Brokers).To(Equal([]BrokerInfo{
				{IsSelf: false, Ordinal: 0, HostName: "broker-0.barco"},
				{IsSelf: false, Ordinal: 6, HostName: "broker-6.barco"},
				{IsSelf: false, Ordinal: 3, HostName: "broker-3.barco"},
				{IsSelf: false, Ordinal: 7, HostName: "broker-7.barco"},
				{IsSelf: false, Ordinal: 1, HostName: "broker-1.barco"},
				{IsSelf: false, Ordinal: 8, HostName: "broker-8.barco"},
				{IsSelf: true, Ordinal: 4, HostName: "broker-4.barco"},
				{IsSelf: false, Ordinal: 9, HostName: "broker-9.barco"},
				{IsSelf: false, Ordinal: 2, HostName: "broker-2.barco"},
				{IsSelf: false, Ordinal: 10, HostName: "broker-10.barco"},
				{IsSelf: false, Ordinal: 5, HostName: "broker-5.barco"},
				{IsSelf: false, Ordinal: 11, HostName: "broker-11.barco"},
			}))
			Expect(topology.LocalIndex).To(Equal(BrokerIndex(6)))
		})
	})

	Describe("Leader()", func() {
		It("should default to the current token when not partition key is provided", func() {
			ordinal := 1
			d := NewDiscoverer(newConfigFake(ordinal), newLocalDbWithNoRecords()).(*discoverer)
			d.k8sClient = &k8sClientFake{6}

			d.Init()
			defer d.Close()

			existingMap := d.generations.Load().(genMap)
			existingMap[d.Topology().MyToken()] = Generation{
				Start:     d.Topology().MyToken(),
				End:       d.Topology().GetToken(d.Topology().LocalIndex + 1),
				Version:   1,
				Leader:    ordinal,
				Followers: []int{4, 2},
				Status:    StatusCommitted,
			}

			info := d.Leader("")
			Expect(info.Leader.Ordinal).To(Equal(ordinal))
			Expect(info.Token).To(Equal(d.Topology().MyToken()))
			Expect(info.Followers[0].Ordinal).To(Equal(4))
			Expect(info.Followers[1].Ordinal).To(Equal(2))
		})

		It("should calculate the primary token and get the generation", func() {
			ordinal := 1
			d := NewDiscoverer(newConfigFake(ordinal), newLocalDbWithNoRecords()).(*discoverer)
			d.k8sClient = &k8sClientFake{6}

			d.Init()
			defer d.Close()

			existingMap := d.generations.Load().(genMap)
			existingMap[d.Topology().GetToken(0)] = Generation{
				Start:     d.Topology().GetToken(0),
				End:       d.Topology().GetToken(1),
				Version:   1,
				Leader:    0,        // Ordinal of node at 0 is zero
				Followers: []int{3}, // Use a single follower as a signal that the generation was obtained
				Status:    StatusCommitted,
			}

			info := d.Leader("a") // token -8839064797231613815 , it should use the first broker
			Expect(info.Leader.Ordinal).To(Equal(0))
			Expect(info.Token).To(Equal(d.Topology().GetToken(0)))
			Expect(info.Followers[0].Ordinal).To(Equal(3))
			// A single follower as a signal
			Expect(info.Followers).To(HaveLen(1))
		})

		It("should set it to the natural owner when there's no information", func() {
			d := NewDiscoverer(newConfigFake(1), newLocalDbWithNoRecords()).(*discoverer)
			d.k8sClient = &k8sClientFake{6}

			d.Init()
			defer d.Close()

			partitionKey := "hui" // token: "7851606034017063987" -> last range
			info := d.Leader(partitionKey)
			Expect(info.Leader.Ordinal).To(Equal(5))
			Expect(info.Token).To(Equal(d.Topology().GetToken(BrokerIndex(5))))
			Expect(info.Followers[0].Ordinal).To(Equal(0))
			Expect(info.Followers[1].Ordinal).To(Equal(3))
		})
	})
})

func newLocalDbWithNoRecords() *dbMocks.Client {
	localDb := new(dbMocks.Client)
	localDb.On("LatestGenerations").Return([]Generation{}, nil)
	return localDb
}

type configFake struct {
	ordinal      int
	baseHostName string
}

func (c *configFake) Ordinal() int {
	return c.ordinal
}

func (c *configFake) BaseHostName() string {
	return c.baseHostName
}

func (c *configFake) ServiceName() string {
	return "barco"
}

func (c *configFake) PodName() string {
	return "barco-0"
}

func (c *configFake) PodNamespace() string {
	return ""
}

func (c *configFake) HomePath() string {
	return "/var/lib/barco"
}

func (c *configFake) ConsumerRanges() int {
	return 8
}

func (c *configFake) ConsumerPort() int {
	return 8081
}

func (c *configFake) ProducerPort() int {
	return 8082
}

func (c *configFake) ClientDiscoveryPort() int {
	return 9089
}

func (c *configFake) ListenOnAllAddresses() bool {
	return true
}

func (c *configFake) DevMode() bool {
	return false
}

func (c *configFake) ShutdownDelay() time.Duration {
	return 0
}

func (c *configFake) FixedTopologyFilePollDelay() time.Duration {
	return 10 * time.Second
}

func newConfigFake(ordinal int) *configFake {
	return &configFake{
		ordinal:      ordinal,
		baseHostName: "barco-",
	}
}

type k8sClientFake struct {
	desiredReplicas int
}

func (c *k8sClientFake) init(_ conf.DiscovererConfig) error {
	return nil
}

func (c *k8sClientFake) getDesiredReplicas() (int, error) {
	return c.desiredReplicas, nil
}

func (c *k8sClientFake) startWatching(replicas int) {

}

func (c *k8sClientFake) replicasChangeChan() <-chan int {
	return make(<-chan int)
}
