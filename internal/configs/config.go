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
	Server               *ServerConfig         `yaml:"server"`
	DataBase             *DataBaseConfig       `yaml:"database"`
	AuthMicroservice     *AuthMicroservice     `yaml:"auth_microservice"`
	CompressMicroservice *CompressMicroservice `yaml:"compress_microservice"`
	NotificationsMicroservice *NotificationsMicroservice `yaml:"notifications_microservice"`
}

// ServerConfig is a struct of server config block in .yaml
type ServerConfig struct {
	Scheme      string `yaml:"scheme"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Front       string `yaml:"frontURI"`
	MediaDir    string `yaml:"mediadir"`
	CVinPDFDir  string `yaml:"cvpdfdir"`
	TamplateDir string `yaml:"tamplateDir"`
	AuthPort    string `yaml:"auth_port"`
	AuthHost    string `yaml:"auth_host"`
	CSRFSecret  string `yaml:"csrf_secret"`
	TLS         *TLS   `yaml:"tls"`
}

type TLS struct {
	Certificate string `yaml:"certificate"`
	Key         string `yaml:"key"`
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

type AuthMicroservice struct {
	Server   *Microservice `yaml:"server"`
	Database *Redis        `yaml:"database"`
}

type CompressMicroservice struct {
	Server             *Microservice `yaml:"server"`
	CompressedMediaDir string        `yaml:"CompressedMediaDir"`
}

type NotificationsMicroservice struct {
	Server   *Microservice `yaml:"server"`
	GRPCserver *Microservice `yaml:"GRPCserver"`
}

type Microservice struct {
	Scheme      string `yaml:"scheme"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	MetricsPort int    `yaml:"metrics_port"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// ReadConfig reads file with configuration.
// Accepts path to .yaml config.
// Returns pointer to Config.
func ReadConfig(confPath string) *Config {
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatalf("unable to read config file: %s", confPath)
	}

	var conf *Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("unable to unmarshal config file")
	}

	if conf.DataBase.ConnectionTimeout == "" {
		conf.DataBase.ConnectionTimeout = "60s"
	}

	return conf
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

func (s *ServerConfig) GetAuthServiceLocation() string {
	return s.AuthHost + ":" + s.AuthPort
}

func (m *Microservice) GetAddress() string {
	return m.Host + ":" + strconv.Itoa(m.Port)
}

func (m *Microservice) GetMetricsAddress() string {
	return m.Host + ":" + strconv.Itoa(m.MetricsPort)
}

func (r *Redis) GetDSN() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}
