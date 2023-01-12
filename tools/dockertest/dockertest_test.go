package dockertest

import (
	"driver/test/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mongoConfig = mock.MongoConfig

func TestDockertest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockertest Suite")
}

var _ = Describe("Dockertest", Ordered, func() {

	Context("NewDockertest", func() {
		It("should be success", func() {
			dockertest := NewDockertest("")
			Expect(dockertest).ShouldNot(BeNil())
		})
	})

	Context("RunMongoDB", func() {
		It("should be success", func() {
			dockertest := NewDockertest("")
			err := dockertest.RunMongoDB(mongoConfig)
			Expect(err).Should(BeNil())

			err = dockertest.Purge()
			Expect(err).Should(BeNil())
		})
	})
})
