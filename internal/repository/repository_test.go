package repository_test

import (
	"driver/internal/repository"
	"driver/pkg"
	"driver/test/mock"
	"driver/tools/dockertest"
	"driver/tools/mongodb"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}

var _ = Describe("Repository", Ordered, func() {

	var mongo mongodb.MongoDBInterface
	var mongoContainer dockertest.DockerTestInterface

	var repo repository.IRepository

	BeforeAll(func() {
		mongoContainer = dockertest.NewDockertest("")
		err := mongoContainer.RunMongoDB(mock.MongoConfig)
		Expect(err).Should(BeNil())

		mongo = mongodb.NewMongoDB(mock.MongoConfig)
		mongo.CreateIndexForGeoJSON()

		mongo.FlushLocations()
	})

	Context("NewRepository", func() {
		It("should be success", func() {
			repo = repository.NewRepository(mongo)
			Expect(repo).ShouldNot(BeNil())
		})
	})

	Context("GetNearestDriver", func() {
		When("there is no driver", func() {
			It("should return error", func() {
				err := mongo.FlushLocations()
				GinkgoWriter.Println(err)

				Expect(err).Should(BeNil())

				res, err := repo.FindNearestDriver(mock.Coordinates)
				GinkgoWriter.Println(fmt.Sprintf("res: %v", res))
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal(pkg.ErrDriverNotFound.Error()))
			})
		})
	})

	Context("BulkCreateDrivers", func() {
		It("should be success", func() {
			err := repo.BulkCreateDrivers(mock.BulkCreateDriversRequest(20))
			Expect(err).Should(BeNil())
		})
	})

	Context("Migration", func() {
		It("should be success", func() {
			err := repo.Migration()
			Expect(err).Should(BeNil())
		})
	})
})
