package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName string = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine the users home directory: %w", err)
	}
	path := filepath.Join(userHome, configFileName)
	return path, nil
}

func Read() (Config,   error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("unable to get the config file path: %w", err)
	}
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open config file: %w", err)
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("unable to parse config file: %w", err)
	}
	return config, nil
}

func (cfg Config) write() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get the config file path: %w", err)
	}
	configFile, err := os.Create(configFilePath)
	defer configFile.Close()
	if err != nil {
		return fmt.Errorf("unable to open config file: %w", err)
	}
	encoder := json.NewEncoder(configFile)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("unable to encode config file: %w", err)
	}
	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return cfg.write()
}
