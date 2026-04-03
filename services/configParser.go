package services

import (
	"encoding/json"
	"os"

	"github.com/tech-thinker/telepath/models"
)

func ParseConfig(configPath string) ([]models.Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config []models.Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
