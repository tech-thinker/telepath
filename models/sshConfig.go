package models

type SSHConfig struct {
	Name     string     `json:"name"`
	Host     string     `json:"host"`
	Port     int        `json:"port"`
	User     string     `json:"user"`
	Password string     `json:"password"`
	KeyFile  string     `json:"keyFile"`
	NextHop  *SSHConfig `json:"nextHop"`
}
