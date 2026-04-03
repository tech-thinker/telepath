package services

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/tech-thinker/telepath/models"
)

func ParseConfig(configPath string) ([]models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Remove single-line comments: // comment
	reSingle := regexp.MustCompile(`(?m)//.*$`)
	data = reSingle.ReplaceAll(data, []byte(""))

	// Remove multi-line comments: /* comment */
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	data = reMulti.ReplaceAll(data, []byte(""))

	// Parse cleaned JSON
	var config []models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}
