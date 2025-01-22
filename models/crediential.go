package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tech-thinker/telepath/constants"
)

type Crediential struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"` // Password or SSH_KEY_FILE
	Password   string `json:"password,omitempty"`
	KeyFile    string `json:"keyFile,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
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
