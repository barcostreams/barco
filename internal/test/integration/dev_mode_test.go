//go:build integration
// +build integration

package integration_test

import (
	"time"

	. "github.com/barcostreams/barco/internal/test/integration"
	. "github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Dev mode", func() {
	var b0 *TestBroker

	AfterEach(func() {
		log.Debug().Msgf("Shutting down dev test cluster")

		if b0 != nil {
			b0.Shutdown()
		}
	})

	It("produces and consumes", func() {
		b0 = NewTestBroker(0, &TestBrokerOptions{DevMode: true})
		b0.WaitOutput("Barco started")

		client := NewTestClient(nil)
		message := `{"hello": "world"}`
		expectOk(client.ProduceJson(0, "abc", message, ""), "should produce json")
		client.RegisterAsConsumer(1, `{"id": "c1", "group": "g1", "topics": ["abc"]}`)

		// Wait for the consumer to be considered
		time.Sleep(1 * time.Second)

		resp := client.ConsumerPoll(0)
		messages := readConsumerResponse(resp)
		expectFindRecord(messages, message)

		time.Sleep(500 * time.Millisecond)
		b0.LookForErrors(30)

	})

	It("supports restarting without cleaning the directory", func() {
		b0 = NewTestBroker(0, &TestBrokerOptions{DevMode: true})
		b0.WaitForStart()
		b0.Shutdown()

		// Restart
		b0.Start()
		b0.WaitForStart()
		time.Sleep(50 * time.Millisecond)

		client := NewTestClient(nil)
		expectOk(client.ProduceJson(0, "abc", `{"hello": "world"}`, ""), "should produce json")
		time.Sleep(200 * time.Millisecond)
		b0.LookForErrors(30)
	})
})
