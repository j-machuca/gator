package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"dbUrl"`
	CurrentUsername string `json:"current_user_name"`
}

func (c *Config) SetUser(current_user_name string) error {
	c.CurrentUsername = current_user_name
	return write(*c)
}

func Read() (Config, error) {
	fp, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(fp)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fp := filepath.Join(homeDir, configFileName)
	return fp, nil
}

func write(c Config) error {
	fp, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}
