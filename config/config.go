package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// EnvType ...
type EnvType string

const (
	// Development environment
	Development EnvType = "development"
	// Test environment
	Test EnvType = "test"
	// Production environment
	Production EnvType = "production"
)

// Option for configurations
type Option struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	HTTP    struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	JwtSecret string `yaml:"jwt_secret"`

	Environment EnvType
}

// Init is using to initialize the configs
func Init(env EnvType) (appConfig *Option, err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	file := path.Join(dir, "../..", "config/config.yaml")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	var options map[EnvType]Option
	err = yaml.Unmarshal(data, &options)
	if err != nil {
		return
	}
	opt := options[env]
	opt.Environment = env
	appConfig = &opt
	return
}
