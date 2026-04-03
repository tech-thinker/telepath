package services

import (
	"github.com/tech-thinker/telepath/models"
)

func ValidateConfig(configs []models.Config) ([]models.Config, error) {
	var validCfgs []models.Config
	for _, cfg := range configs {
		err := cfg.Validate()
		if err != nil {
			return nil, err
		}
		validCfgs = append(validCfgs, cfg)
	}
	return validCfgs, nil
}
