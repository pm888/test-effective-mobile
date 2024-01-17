package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	ServerHost string
	ServerPort string
	DbHost     string
}

func New() (*Config, error) {
	var c = new(Config)
	if err := c.LoadEnvFile(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) LoadEnvFile() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	c.DbPort = getEnv("DB_PORT")
	c.DbUser = getEnv("DB_USER")
	c.DbPassword = getEnv("DB_PASSWORD")
	c.DbName = getEnv("DB_NAME")
	c.DbHost = getEnv("SERVER_HOST")
	c.ServerPort = getEnv("SERVER_PORT")
	c.ServerHost = getEnv("DB_HOST")
	return nil
}
func getEnv(key string) string {
	return os.Getenv(key)
}
