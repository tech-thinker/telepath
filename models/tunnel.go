package models

import (
	"encoding/json"
	"errors"
)

type Tunnel struct {
	Name       string   `json:"name,omitempty"`
	LocalPort  int      `json:"localPort,omitempty"`
	RemoteHost string   `json:"remoteHost,omitempty"`
	RemotePort int      `json:"remotePort,omitempty"`
	HostChain  []string `json:"hostChain,omitempty"`
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
