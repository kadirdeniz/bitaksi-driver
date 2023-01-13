package pact

import (
	"driver/tools/fiber"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"testing"
)

var testingObject *testing.T

func TestProvier(t *testing.T) {
	testingObject = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Provier Suite Test")
}

var _ = Describe("Provier Suite", Ordered, func() {

	GinkgoWriter.Println("Starting Pact Test")
	var pact *dsl.Pact
	var verifier types.VerifyRequest

	port, _ := utils.GetFreePort()
	const host = "127.0.0.1"

	BeforeAll(func() {
		pact = &dsl.Pact{
			Consumer: "MatcherService",
			Provider: "DriverService",
			Host:     host,
			LogDir:   "logs",
			LogLevel: "DEBUG",
		}

		verifier = types.VerifyRequest{
			ProviderBaseURL:            fmt.Sprintf("http://%s:%d", host, port),
			ProviderVersion:            "1.0.0",
			Tags:                       []string{"main"},
			PactURLs:                   []string{"test/pact/pacts/matcherservice-driverservice.json"},
			PublishVerificationResults: false,
			FailIfNoPactsFound:         true,
		}

		go fiber.StartServer(port)

	})

	Context("Provider", func() {
		It("should be ok", func() {
			_, err := pact.VerifyProvider(testingObject, verifier)
			Expect(err).Should(BeNil())
		})
	})
})
