// Package configs contains configs of project
package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Config is a struct of .yaml config file
type Config struct {
	Server   *ServerConfig   `yaml:"server"`
	DataBase *DataBaseConfig `yaml:"database"`
}

// ServerConfig is a struct of server config block in .yaml
type ServerConfig struct {
	Scheme   string `yaml:"scheme"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Front    string `yaml:"frontURI"`
	mediaDir string `yaml:"mediadir"`
}

type DataBaseConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	Schema            string `yaml:"schema"`
	DBName            string `yaml:"db_name"`
	SSLMode           string `yaml:"ssl_mode"`
	ConnectionTimeout string `yaml:"conn_timeout,omitempty"` // Temporary unused
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

	if conf.DataBase.ConnectionTimeout == "" {
		conf.DataBase.ConnectionTimeout = "60s"
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

// GetFrontURI returns front uri. E.g http://127.0.0.1:3000
func (s *ServerConfig) GetFrontURI() string {
	return s.Front
}

// GetAddress returns directory with user's files
func (s *ServerConfig) GetMediaDir() string {
	return s.mediaDir
}

func (d *DataBaseConfig) GetDSN() string {
	return fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
		d.User,
		d.DBName,
		d.Password,
		d.Host,
		d.Port,
		d.SSLMode,
	)
}
