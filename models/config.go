package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tech-thinker/telepath/constants"
)

type Config struct {
	Hosts        map[string]HostConfig  `json:"hosts"`
	Credientials map[string]Crediential `json:"credientials"`
	Tunnels      map[string]Tunnel      `json:"tunnels"`
}

type Tunnel struct {
	Name       string   `json:"name"`
	LocalPort  int      `json:"localPort"`
	RemoteHost string   `json:"remoteHost"`
	RemotePort int      `json:"remotePort"`
	HostChain  []string `json:"hostChain"`
}

func (t *Tunnel) ToByte() []byte {
	data, _ := json.Marshal(t)
	return data
}

func (t *Tunnel) Validate(hosts map[string]HostConfig) (bool, error) {
	if len(t.Name) == 0 {
		return false, errors.New("Tunnel name must not be empty.")
	}

	if t.LocalPort <= 0 {
		return false, errors.New("Local port invalid.")
	}

	if len(t.RemoteHost) == 0 {
		return false, errors.New("Remote host must not be empty.")
	}

	if t.RemotePort <= 0 {
		return false, errors.New("Remote port invalid.")
	}

	if len(t.HostChain) == 0 {
		return false, errors.New("Host chain must not be empty.")
	}

	for _, host := range t.HostChain {
		_, ok := hosts[host]
		if !ok {
			return false, errors.New("Host doesn't exists.")
		}
	}

	return true, nil
}

type HostConfig struct {
	Name            string `json:"name"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	CredientialName string `json:"credientialName"`
}

func (h *HostConfig) ToByte() []byte {
	data, _ := json.Marshal(h)
	return data
}

func (h *HostConfig) Validate(creds map[string]Crediential) (bool, error) {
	if len(h.Name) == 0 {
		return false, errors.New("Host name must not be empty.")
	}

	if len(h.Host) == 0 {
		return false, errors.New("Host address must not be empty.")
	}

	if h.Port <= 0 {
		return false, errors.New("Host port invalid.")
	}

	if len(h.User) == 0 {
		return false, errors.New("Host user must not be empty.")
	}

	if len(h.CredientialName) == 0 {
		return false, errors.New("Host crediential name must not be empty.")
	}

	_, ok := creds[h.CredientialName]
	if !ok {
		return false, errors.New("Host crediential doesn't exists.")
	}

	return true, nil
}

type Crediential struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // Password or SSH_KEY_FILE
	Password string `json:"password"`
	KeyFile  string `json:"keyFile"`
}

func (c *Crediential) ToByte() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *Crediential) Validate() (bool, error) {
	if len(c.Name) == 0 {
		return false, errors.New("Crediential name must not be empty.")
	}

	if c.Type != constants.CREDIENTIAL_PASS && c.Type != constants.CREDIENTIAL_KEY {
		return false, errors.New("Crediential type must be PASS or KEY.")
	}

	if c.Type == constants.CREDIENTIAL_PASS && len(c.Password) == 0 {
		return false, errors.New("Crediential password must not be empty.")
	}

	if c.Type == constants.CREDIENTIAL_KEY && len(c.KeyFile) == 0 {
		return false, errors.New("Crediential key file path must not be empty.")
	}

	if c.Type == constants.CREDIENTIAL_KEY && len(c.KeyFile) > 0 {
		if _, err := os.Stat(c.KeyFile); os.IsNotExist(err) {
			return false, fmt.Errorf("crediential key file '%s' does not exist", c.KeyFile)
		}
	}
	return true, nil
}
