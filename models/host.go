package models

import (
	"encoding/json"
	"errors"
)

type HostConfig struct {
	Name            string `json:"name,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	User            string `json:"user,omitempty"`
	CredientialName string `json:"credientialName,omitempty"`
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
