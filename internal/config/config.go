package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	file, err := openConfigFile(os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Config) SetUser(userName string) error {
	config.CurrentUserName = userName

	file, err := openConfigFile(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

func openConfigFile(flags int) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := homeDir + "/" + configFileName
	file, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
