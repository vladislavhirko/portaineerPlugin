package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	API API `toml:"api"`
	LevelDB Level      `toml:"level"`
	MClient Mattermost `toml:"mattermost"`
	PClient Portainer  `toml:"portainer"`
}

type API struct{
	Port string `toml:"port"`
}

type Level struct {
	Path string `toml:"path"`
}

type Mattermost struct {
	Address  string `toml:"address"`
	Port     string `toml:"port"`
	Email    string `toml:"email"`
	Password string `toml:"password"`
}

type Portainer struct {
	Login         string `toml:"login"`
	Password      string `toml:"password"`
	Address       string `toml:"address"`
	Port          string `toml:"port"`
	CheckInterval string `toml:"check_interval"`
}

func GetConfig() Config {
	config := Config{
		API: *new(API),
		LevelDB: *new(Level),
		MClient: *new(Mattermost),
		PClient: *new(Portainer),
	}
	_, err := toml.DecodeFile("config/config.toml", &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}
