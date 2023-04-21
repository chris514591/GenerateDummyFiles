package Fileutil

import (
	"encoding/json"
	"os"
)

type Config struct {
	Path          string `json:"path"`
	MinNumOfFiles int    `json:"min_num_of_files"`
	MaxNumOfFiles int    `json:"max_num_of_files"`
	MinFileSize   int    `json:"min_file_size"`
	MaxFileSize   int    `json:"max_file_size"`
}

func ReadConfigFile(filename string) (Config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func OpenErrorsLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}
