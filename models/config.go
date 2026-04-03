package models

import (
	"fmt"
	"slices"

	"github.com/tech-thinker/telepath/constants"
)

type Server struct {
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	Username   string  `json:"username"`
	AuthType   string  `json:"authType"`
	Password   string  `json:"password"`
	Key        string  `json:"key"`
	Passphrase string  `json:"passphrase"`
	Jump       *Server `json:"jump"`
}

func (s *Server) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *Server) Validate(fieldName string) error {
	if len(s.Host) == 0 {
		return fmt.Errorf("`%s.host` must be present.", fieldName)
	}

	if s.Port == 0 {
		return fmt.Errorf("`%s.port` is invalid.", fieldName)
	}

	if len(s.Username) == 0 {
		return fmt.Errorf("`%s.username` must be present.", fieldName)
	}

	if !slices.Contains([]string{constants.CREDIENTIAL_KEY, constants.CREDIENTIAL_PASS}, s.AuthType) {
		return fmt.Errorf("`%s.authType` must be present.", fieldName)
	}

	if s.AuthType == constants.CREDIENTIAL_KEY && len(s.Key) == 0 {
		return fmt.Errorf("`%s.key` must be present for authType KEY.", fieldName)
	}

	if s.AuthType == constants.CREDIENTIAL_PASS && len(s.Password) == 0 {
		return fmt.Errorf("`%s.password` must be present for authType PASS.", fieldName)
	}

	if s.Jump != nil {
		err := s.Jump.Validate(fmt.Sprintf("%s.jump", fieldName))
		if err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	LocalPort  int     `json:"localPort"`
	LocalHost  string  `json:"localHost"`
	RemotePort int     `json:"remotePort"`
	RemoteHost string  `json:"remoteHost"`
	Server     *Server `json:"server"`
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

func (c *Config) Validate() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("`name` must be present.")
	}

	if !slices.Contains([]string{"L", "R"}, c.Type) {
		return fmt.Errorf("`type` must be L or R.")
	}

	if c.LocalPort <= 0 {
		return fmt.Errorf("`localPort` is invalid.")
	}

	if len(c.LocalHost) == 0 {
		c.LocalHost = "0.0.0.0"
	}

	if c.RemotePort <= 0 {
		return fmt.Errorf("`remotePort` is invalid.")
	}

	if len(c.RemoteHost) == 0 {
		c.RemoteHost = "0.0.0.0"
	}

	if c.Server == nil {
		return fmt.Errorf("`server` must be present.")
	}

	err := c.Server.Validate("server")
	if err != nil {
		return err
	}
	return nil
}
