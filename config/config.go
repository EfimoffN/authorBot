package config

import (
	"os"

	"github.com/EfimoffN/authorBot/config/lib/e"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken        string `yaml: "bot_token"`
	BindAddr        string `yaml: "bind_addr"`
	LogLevel        string `yaml: "log_level"`
	ConnectPostgres string `yaml: "connect_postgres"`
}

func CreateConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, e.Wrap("open file", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, e.Wrap("decoder file", err)
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		t := "stat path: '%s'"
		return e.Wrap("stat path: '%s'", err)
	}
}

func validateConfig(config Config) error {

	return nil
}
