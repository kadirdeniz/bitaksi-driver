package pkg

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var AppConfigs Configs

type MongoDBConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type Application struct {
	API_KEY string `yaml:"api_key"`
	PORT    string `yaml:"port"`
}

type Configs struct {
	MongoDB     MongoDBConfig `yaml:"mongodb"`
	Application Application   `yaml:"application"`
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "./..")

	err := readYamlFile(filepath.Join(basepath, "./configs/mongodb.yaml"))
	if err != nil {
		panic(err)
	}

	err = readYamlFile(filepath.Join(basepath, "./configs/application.yaml"))
	if err != nil {
		panic(err)
	}
}

func readYamlFile(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &AppConfigs)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configs) GetMongoDBConfig() MongoDBConfig {
	return c.MongoDB
}

func (c *Configs) GetApplicationConfig() Application {
	return c.Application
}
