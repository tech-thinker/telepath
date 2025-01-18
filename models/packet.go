package models

import (
	"encoding/json"
)

type Packet struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Data   []byte `json:"data"`
}

func (p *Packet) ToByte() []byte {
	data, _ := json.Marshal(p)
	return data
}
