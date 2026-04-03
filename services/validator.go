package services

import (
	"fmt"

	"github.com/tech-thinker/telepath/models"
)

func ValidateConfig(configs []models.Config) ([]models.Config, error) {
	for _, cfg := range configs {
		if cfg.Type != "L" && cfg.Type != "R" {
			return nil, fmt.Errorf("invalid tunnel type %s for %s", cfg.Type, cfg.Name)
		}
	}
	return configs, nil
}
