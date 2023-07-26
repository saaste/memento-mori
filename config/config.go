package config

import (
	"encoding/json"
	"os"

	"github.com/saaste/memento-mori/utils"
)

type AppConfig struct {
	Birthday       utils.Date
	LifeExpectancy int
	Events         []Event
}

type Event struct {
	Date  utils.Date
	Title string
	Label string
}

func ReadConfig() (AppConfig, error) {
	var config AppConfig

	f, err := os.ReadFile("./config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal([]byte(f), &config)
	if err != nil {
		return config, err
	}
	return config, nil

}
