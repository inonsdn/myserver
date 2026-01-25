package config

import "fmt"

type Config struct {
	host string
	port int
}

type ConfigFunc func(c *Config)

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func defaultConfig() *Config {
	return &Config{
		host: "localhost",
		port: 8080,
	}
}

func LoadConfig(cf ...ConfigFunc) *Config {
	config := defaultConfig()
	for _, f := range cf {
		f(config)
	}
	return config
}
