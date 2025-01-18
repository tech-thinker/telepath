package utils

import (
	"encoding/json"
	"log"

	"github.com/tech-thinker/telepath/models"
)

func ToPacket(buff []byte) (models.Packet, error) {
	var packet models.Packet
	err := json.Unmarshal(buff, &packet)
	if err != nil {
		log.Printf("Error decoding packet: %v", err)
		return packet, err
	}
	return packet, nil
}

func ToCrediential(buff []byte) (models.Crediential, error) {
	var cred models.Crediential
	err := json.Unmarshal(buff, &cred)
	if err != nil {
		return cred, err
	}
	return cred, nil
}

func ToHostConfig(buff []byte) (models.HostConfig, error) {
	var hostConfig models.HostConfig
	err := json.Unmarshal(buff, &hostConfig)
	if err != nil {
		return hostConfig, err
	}
	return hostConfig, nil
}

func ToTunnel(buff []byte) (models.Tunnel, error) {
	var tunnel models.Tunnel
	err := json.Unmarshal(buff, &tunnel)
	if err != nil {
		return tunnel, err
	}
	return tunnel, nil
}
