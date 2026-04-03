package models

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

type Config struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	LocalPort  int     `json:"localPort"`
	LocalHost  string  `json:"localHost"`
	RemotePort int     `json:"remotePort"`
	RemoteHost string  `json:"remoteHost"`
	Server     *Server `json:"server"`
}
