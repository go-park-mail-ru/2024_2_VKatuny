// Package configs contains configs of project
package configs

import (
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Config is a struct of .yaml config file
type Config struct {
	Server *ServerConfig `yaml:"server"`
}

// ServerConfig is a struct of server config block in .yaml
type ServerConfig struct {
	Scheme             string `yaml:"scheme"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Front              string `yaml:"frontURI"`
	MediaDir           string `yaml:"mediadir"`
	CompressedMediaDir string `yaml:"CompressedMediaDir"`
}

// ReadConfig reads file with configuration.
// Accepts path to .yaml config.
// Returns pointer to Config.
func ReadConfig(confPath string) (*Config, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatalf("unable to read config file: %s", confPath)
		return nil, err
	}

	var conf *Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("unable to unmarshal config file")
	}

	return conf, nil
}

// GetAddress returns address of server
func (s *ServerConfig) GetAddress() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// GetHostWithScheme returns host with scheme and port
// e.g. http://127.0.0.1:8080
func (s *ServerConfig) GetHostWithScheme() string {
	return s.Scheme + "://" + s.Host
}
