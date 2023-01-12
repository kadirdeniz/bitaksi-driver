package mongodb_test

import (
	"driver/pkg"
	"driver/test/mock"
	"driver/tools/dockertest"
	"driver/tools/mongodb"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mongoConfig = mock.MongoConfig

func TestMongodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}

var _ = Describe("MongoDB", Ordered, func() {

	var mongo mongodb.MongoDBInterface
	var dockerContainer *dockertest.Dockertest

	BeforeAll(func() {
		dockerContainer = dockertest.NewDockertest("")
		err := dockerContainer.RunMongoDB(mongoConfig)
		Expect(err).Should(BeNil())

		mongodb.NewMongoDB(mongoConfig).FlushLocations()
	})

	AfterAll(func() {
		dockerContainer.Purge()
	})

	Context("Connect", func() {
		It("Should return database", func() {
			mongo = mongodb.NewMongoDB(mongoConfig)
			db, err := mongo.Connect()
			Expect(err).Should(BeNil())
			Expect(db).ShouldNot(BeNil())
		})
	})

	Context("GetMongoDBURI", func() {
		It("should return mongodb uri", func() {
			mongodbURI := mongo.GetMongoDBURI()
			Expect(mongodbURI).Should(Equal("mongodb://" + mongoConfig.Username + ":" + mongoConfig.Password + "@" + mongoConfig.Host + ":" + mongoConfig.Port))
		})
	})

	Context("GetUserCollection", func() {
		It("should return mongodb user collection", func() {
			userCollection := mongo.GetDriverCollection()
			Expect(userCollection).ShouldNot(BeNil())
		})
	})

	Context("GetDatabase", func() {
		It("should return mongodb database", func() {
			database := mongo.GetDatabase()
			Expect(database).ShouldNot(BeNil())
		})
	})

	Context("BulkCreateDrivers", func() {
		It("should be success", func() {
			err := mongo.BulkCreateDrivers(mock.BulkCreateDriversRequest(20))
			Expect(err).Should(BeNil())

			drivers, err := mongo.FindLocations()
			Expect(err).Should(BeNil())
			Expect(len(drivers)).Should(Equal(20))
		})
	})

	Context("Flush Locations", func() {
		It("should delete locations", func() {
			err := mongo.FlushLocations()
			Expect(err).Should(BeNil())
		})

		It("shouldnt return locations", func() {
			users, err := mongo.FindLocations()
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(Equal(pkg.ErrDriverNotFound))
			Expect(users).Should(BeNil())
		})
	})
})
