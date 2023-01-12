package mongodb

import (
	"driver/test/mock"
	"driver/tools/dockertest"
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

	var mongo MongoDBInterface
	var dockerContainer *dockertest.Dockertest

	BeforeAll(func() {
		dockerContainer = dockertest.NewDockertest("")
		err := dockerContainer.RunMongoDB(mongoConfig)
		Expect(err).Should(BeNil())

		NewMongoDB(mongoConfig).FlushUsers()
	})

	AfterAll(func() {
		dockerContainer.Purge()
	})

	Context("Connect", func() {
		It("Should return database", func() {
			mongo = NewMongoDB(mongoConfig)
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

})
