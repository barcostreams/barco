package interbroker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	cMocks "github.com/barcostreams/barco/internal/test/conf/mocks"
	dMocks "github.com/barcostreams/barco/internal/test/discovery/mocks"
	. "github.com/barcostreams/barco/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Gossiper", func() {
	Describe("ReadProducerOffset", func() {
		It("Should call the server", func() {
			var mu sync.Mutex
			var requestUrl *url.URL
			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				requestUrl = r.URL
				defer mu.Unlock()
				fmt.Fprintln(w, "123")
			}))
			ts.EnableHTTP2 = true
			ts.Start()
			defer ts.Close()

			port, err := strconv.Atoi(strings.Split(ts.URL, ":")[2])
			Expect(err).NotTo(HaveOccurred())

			const ordinal = 0
			topicId := TopicDataId{
				Name:       "abc",
				Token:      1,
				RangeIndex: 2,
				Version:    3,
			}
			discoverer := new(dMocks.Discoverer)
			discoverer.On("CurrentOrPastBroker", mock.Anything).Return(&BrokerInfo{HostName: "127.0.0.1"})
			config := new(cMocks.Config)
			config.On("GossipPort").Return(port)

			g := gossiper{
				connections: atomic.Value{},
				discoverer:  discoverer,
				config:      config,
			}

			clients := make(clientMap)
			clients[ordinal] = &clientInfo{
				gossipClient: ts.Client(),
				isConnected:  1,
			}
			g.connections.Store(clients)

			value, err := g.ReadProducerOffset(ordinal, &topicId)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(int64(123)))
			mu.Lock()
			Expect(requestUrl).NotTo(BeNil())
			Expect(requestUrl.Path).To(Equal("/v1/producer/offset/abc/1/2/3"))
			mu.Unlock()
		})
	})
})

func newTestTopology(length int, ordinal int) *TopologyInfo {
	brokers := make([]BrokerInfo, length)
	for i := 0; i < length; i++ {
		brokers[i] = BrokerInfo{
			IsSelf:   i == ordinal,
			Ordinal:  i,
			HostName: fmt.Sprintf("127.0.0.%d", i+1),
		}
	}

	t := NewTopology(brokers, ordinal)
	return &t
}
