package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
)

type Configuration interface {
	LoadConfig()
	SaveConfig() error
	Config() *models.Config
}

type configuration struct {
	cfgDir      string
	cfgFile     string
	cfgFilePath string
	config      *models.Config
}

func (cfg *configuration) Config() *models.Config {
	return cfg.config
}

func (cfg *configuration) isConfigExists() {
	if _, err := os.Stat(cfg.cfgDir); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.cfgDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat(cfg.cfgFilePath); os.IsNotExist(err) {
		cfg.config = &models.Config{
			Hosts:        make(map[string]models.HostConfig),
			Credientials: make(map[string]models.Crediential),
			Tunnels:      make(map[string]models.Tunnel),
		}
		err = cfg.SaveConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cfg *configuration) LoadConfig() {
	data, err := os.ReadFile(cfg.cfgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, cfg.config)
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *configuration) SaveConfig() error {
	data, err := json.MarshalIndent(cfg.config, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(cfg.cfgFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig() Configuration {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory: ", err)
	}
	cfgDir := filepath.Join(homeDir, constants.CONFIG_DIR)
	cfgFile := constants.CONFIG_FILE
	cfgFilePath := filepath.Join(cfgDir, cfgFile)
	cfg := configuration{
		cfgDir:      cfgDir,
		cfgFile:     cfgFile,
		cfgFilePath: cfgFilePath,
		config: &models.Config{
			Hosts:        make(map[string]models.HostConfig),
			Credientials: make(map[string]models.Crediential),
			Tunnels:      make(map[string]models.Tunnel),
		},
	}
	cfg.isConfigExists()
	cfg.LoadConfig()
	return &cfg
}
