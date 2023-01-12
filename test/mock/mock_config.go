package mock

import "driver/pkg"

var MongoConfig = pkg.MongoDBConfig{
	Username: "admin",
	Password: "admin",
	Host:     "localhost",
	Port:     "27017",
	Database: "bitaksi",
}
