package configs

import (
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Struct of .yaml config
type Config struct {
	Server *ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Scheme string `yaml:"scheme"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

// Reads file with configuration.
// Accepts path to .yaml config.
// Returns pointer to Config.
func ReadConfig(confPath string) (*Config, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal("unable to read config file")
		return nil, err
	}

	var conf *Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("unable to unmarshal config file")
	}

	return conf, nil
}

// Returns address of server
func (s *ServerConfig) GetAddress() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

func (s *ServerConfig) GetHostWithScheme() string {
	return s.Scheme + "://" + s.Host
}
