package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Telegram Telegram `json:"telegram'`
}

type Telegram struct {
	URI    string `json:"uri"`
	ChatID string `json:"chat_id"`
	Token  string `json:"token"`
}

func LoadConfig() error {
	return config()
}

func config() error {
	conf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("load config, read file: %w", err)
	}

	return readJsonfile(conf)
}

func readJsonfile(data []byte) error {
	c := Config{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return fmt.Errorf("unmarshal json config file: %w", err)
	}
	return setEnv(configToMap(c))
}

func configToMap(c Config) map[string]string {
	return map[string]string{
		"URI":     c.Telegram.URI,
		"TOKEN":   c.Telegram.Token,
		"CHAT_ID": c.Telegram.ChatID,
	}
}

func setEnv(envs map[string]string) error {
	for key, value := range envs {
		if value == "" {
			return fmt.Errorf("set empty env :key %s has no value", key)
		}
		err := os.Setenv(fmt.Sprintf("TELEGRAME_%s", key), value)
		if err != nil {
			return fmt.Errorf("set env key %s value %s: %w", key, value, err)
		}
	}
	return nil
}
