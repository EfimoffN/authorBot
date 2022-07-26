package config

import (
	"fmt"
	"os"

	"github.com/EfimoffN/authorBot/lib/e"
	"gopkg.in/yaml.v3"
)

type ConfigApp struct {
	BotToken        string `yaml: "bot_token"`
	BindAddr        string `yaml: "bind_addr"`
	LogLevel        string `yaml: "log_level"`
	ConnectPostgres string `yaml: "connect_postgres"`
	Timeout         int    `yaml: "timeout"`
}

func CreateConfig(configPath string) (*ConfigApp, error) {

	err := validateConfigAppPath(configPath)
	if err != nil {
		return nil, e.Wrap("validat config path", err)
	}

	config := &ConfigApp{}

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

func validateConfigAppPath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return e.Wrap(fmt.Sprintf("config path: '%s'", path), err)
	}

	return nil
}

// func validateConfig(config ConfigApp) error {

// 	return nil
// }
