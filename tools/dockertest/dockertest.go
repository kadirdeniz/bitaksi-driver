package dockertest

import (
	"driver/pkg"
	"driver/tools/mongodb"
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"log"
)

type DockerTest struct{}

var MongoDBEnvirontmentVariables = []string{
	"MONGO_INITDB_ROOT_USERNAME=" + pkg.AppConfigs.MongoDB.Username,
	"MONGO_INITDB_ROOT_PASSWORD=" + pkg.AppConfigs.MongoDB.Password,
	"MONGO_INITDB_DATABASE=" + pkg.AppConfigs.MongoDB.Database,
}

const MongoDBImage = "mongo"
const MongoDBTag = "5.0"
const ExposedPort = "27017"

type Dockertest struct {
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

type DockerTestInterface interface {
	RunMongoDB(config pkg.MongoDBConfig) error
	Purge() error
}

func NewDockertest(endpoint string) DockerTestInterface {
	pool, err := dockertest.NewPool(endpoint)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return &Dockertest{
		Pool: pool,
	}
}

func (d *Dockertest) RunMongoDB(config pkg.MongoDBConfig) error {

	var err error

	d.Resource, err = d.Pool.Run(MongoDBImage, MongoDBTag, MongoDBEnvirontmentVariables)
	if err != nil {
		return errors.New("Could not start resource")
	}

	if err = d.Pool.Retry(func() error {
		_, err := mongodb.NewMongoDB(config).Connect()
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}); err != nil {
		return errors.New("Could not connect to docker")
	}

	return nil
}

func (d *Dockertest) Purge() error {
	if err := d.Pool.Purge(d.Resource); err != nil {
		return errors.New("Could not purge resource")
	}
	return nil
}
