package constants

const (
	CONFIG_DIR  = ".config/telepath"
	CONFIG_FILE = "telepath.json"
)

const (
	PID_FILE_PATH = "/tmp/telepath.pid"
	SOCKET_PATH   = "/tmp/telepath.sock"
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
