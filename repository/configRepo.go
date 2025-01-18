package repository

import (
	"errors"
	"strings"

	"github.com/tech-thinker/telepath/config"
	"github.com/tech-thinker/telepath/models"
)

type ConfigRepo interface {
	AddCrediential(name string, t string, password string, keyFile string) error
	RemoveCrediential(name string) error
	DetailCrediential(name string) (models.Crediential, error)
	ListCrediential() ([]models.Crediential, error)
	AddHost(name string, host string, port int, user string, credientialName string) error
	RemoveHost(name string) error
	DetailHost(name string) (models.HostConfig, error)
	ListHost() ([]models.HostConfig, error)
	AddTunnel(name string, localPort int, remoteHost string, remotePort int, hostChain []string) error
	RemoveTunnel(name string) error
	DetailTunnel(name string) (models.Tunnel, error)
	ListTunnel() ([]models.Tunnel, error)
}

type configRepo struct {
	cfg config.Configuration
}

func (repo *configRepo) AddCrediential(name string, t string, password string, keyFile string) error {
	cred := models.Crediential{
		Name:     strings.TrimSpace(name),
		Type:     strings.TrimSpace(t),
		Password: strings.TrimSpace(password),
		KeyFile:  strings.TrimSpace(keyFile),
	}
	ok, err := cred.Validate()
	if !ok {
		return err
	}

	// Add Crediential
	repo.cfg.Config().Credientials[cred.Name] = cred
	err = repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) RemoveCrediential(name string) error {
	delete(repo.cfg.Config().Credientials, name)
	err := repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) DetailCrediential(name string) (models.Crediential, error) {
	cred, ok := repo.cfg.Config().Credientials[name]
	if !ok {
		return cred, errors.New(`Crediential not found.`)
	}
	return cred, nil
}

func (repo *configRepo) ListCrediential() ([]models.Crediential, error) {
	list := []models.Crediential{}
	for _, val := range repo.cfg.Config().Credientials {
		list = append(list, val)
	}
	return list, nil
}

func (repo *configRepo) AddHost(name string, host string, port int, user string, credientialName string) error {
	hostCfg := models.HostConfig{
		Name:            strings.TrimSpace(name),
		Host:            strings.TrimSpace(host),
		Port:            port,
		User:            strings.TrimSpace(user),
		CredientialName: strings.TrimSpace(credientialName),
	}
	ok, err := hostCfg.Validate(repo.cfg.Config().Credientials)
	if !ok {
		return err
	}

	// Add Host
	repo.cfg.Config().Hosts[hostCfg.Name] = hostCfg
	err = repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) RemoveHost(name string) error {
	delete(repo.cfg.Config().Hosts, name)
	err := repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) DetailHost(name string) (models.HostConfig, error) {
	host, ok := repo.cfg.Config().Hosts[name]
	if !ok {
		return host, errors.New(`Host not found.`)
	}
	return host, nil
}

func (repo *configRepo) ListHost() ([]models.HostConfig, error) {
	list := []models.HostConfig{}
	for _, val := range repo.cfg.Config().Hosts {
		list = append(list, val)
	}
	return list, nil
}

func (repo *configRepo) AddTunnel(name string, localPort int, remoteHost string, remotePort int, hostChain []string) error {
	tunnel := models.Tunnel{
		Name:       strings.TrimSpace(name),
		LocalPort:  localPort,
		RemoteHost: strings.TrimSpace(remoteHost),
		RemotePort: remotePort,
		HostChain:  hostChain,
	}
	ok, err := tunnel.Validate(repo.cfg.Config().Hosts)
	if !ok {
		return err
	}

	// Add Host
	repo.cfg.Config().Tunnels[tunnel.Name] = tunnel
	err = repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) RemoveTunnel(name string) error {
	delete(repo.cfg.Config().Tunnels, name)
	err := repo.cfg.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (repo *configRepo) DetailTunnel(name string) (models.Tunnel, error) {
	tunnel, ok := repo.cfg.Config().Tunnels[name]
	if !ok {
		return tunnel, errors.New(`Tunnel not found.`)
	}
	return tunnel, nil
}

func (repo *configRepo) ListTunnel() ([]models.Tunnel, error) {
	list := []models.Tunnel{}
	for _, val := range repo.cfg.Config().Tunnels {
		list = append(list, val)
	}
	return list, nil
}

func NewConfigRepo(cfg config.Configuration) ConfigRepo {
	return &configRepo{
		cfg: cfg,
	}
}
