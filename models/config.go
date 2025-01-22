package models

type Config struct {
	Hosts        map[string]HostConfig  `json:"hosts,omitempty"`
	Credientials map[string]Crediential `json:"credientials,omitempty"`
	Tunnels      map[string]Tunnel      `json:"tunnels,omitempty"`
}
