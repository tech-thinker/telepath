package constants

const (
	CONFIG_DIR  = ".config/telepath"
	CONFIG_FILE = "telepath.json"
	PID_FILE    = "telepath.pid"
	SOCKET_FILE = "telepath.sock"
	TCP_ADDR    = "127.0.0.1:54321"
)

const (
	CREDIENTIAL_PASS = "PASS"
	CREDIENTIAL_KEY  = "KEY"
)

const (
	PACKET_TYPE_HOST_CONFIG = "HostConfig"
	PACKET_TYPE_CREDIENTIAL = "Crediential"
	PACKET_TYPE_TUNNEL      = "Tunnel"
	PACKET_TYPE_STRING      = "string"
)
