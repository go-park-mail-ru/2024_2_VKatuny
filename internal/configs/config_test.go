package configs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerGetAddress(t *testing.T) {
	serverConfig := ServerConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
	require.Equal(t, "127.0.0.1:8080", serverConfig.GetAddress())
}

func TestServerGetHostWithScheme(t *testing.T) {
	serverConfig := ServerConfig{
		Scheme: "http",
		Host:   "127.0.0.1",
		Port:   8080,
	}
	require.Equal(t, "http://127.0.0.1", serverConfig.GetHostWithScheme())
}

func TestDBGetDSN(t *testing.T) {
	dbConfig := DataBaseConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
	expected := "user=postgres dbname=postgres password=postgres host=127.0.0.1 port=5432 sslmode=disable"
	require.Equal(t, expected, dbConfig.GetDSN())
}

func TestMicroserviceGetAddress(t *testing.T) {
	microserviceConfig := Microservice{
		Host: "127.0.0.1",
		Port: 8080,
	}
	require.Equal(t, "127.0.0.1:8080", microserviceConfig.GetAddress())
}

func TestRedisGetDSN(t *testing.T) {
	redisConfig := Redis{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "password",
	}
	require.Equal(t, "127.0.0.1:6379", redisConfig.GetDSN())
}

func TestGetMetricsAddress(t *testing.T) {
	microserviceConfig := Microservice{
		Host: "127.0.0.1",
		Port: 8080,
	}
	require.Equal(t, "127.0.0.1:0", microserviceConfig.GetMetricsAddress())
}

func TestGetAuthServiceLocation(t *testing.T) {
	serviceConfig := ServerConfig{
		AuthHost: "127.0.0.1",
		AuthPort: "8080",
	}
	require.Equal(t, "127.0.0.1:8080", serviceConfig.GetAuthServiceLocation())
}
