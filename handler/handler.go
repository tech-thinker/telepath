package handler

import (
	"context"

	"github.com/tech-thinker/telepath/config"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/repository"
)

type Handler interface {
	Handle(ctx context.Context, packet models.Packet) (models.Packet, error)
}

type handler struct {
	configRepo repository.ConfigRepo
	socketRepo repository.SocketRepo
}

func (h *handler) Handle(ctx context.Context, packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	switch packet.Action {
	case "add-cred":
		return h.addCrediential(packet)
	case "remove-cred":
		return h.removeCrediential(packet)
	case "detail-cred":
		return h.detailCrediential(packet)
	case "list-cred":
		return h.listCrediential(packet)

	case "add-host":
		return h.addHost(packet)
	case "remove-host":
		return h.removeHost(packet)
	case "detail-host":
		return h.detailHost(packet)
	case "list-host":
		return h.listHost(packet)

	case "add-tunnel":
		return h.addTunnel(packet)
	case "remove-tunnel":
		return h.removeTunnel(packet)
	case "detail-tunnel":
		return h.detailTunnel(packet)
	case "list-tunnel":
		return h.listTunnel(packet)

	case "start-tunnel":
		return h.startTunnel(packet)
	case "stop-tunnel":
		return h.stopTunnel(packet)
	}

	return result, nil
}

func NewHandler(cfg config.Configuration) Handler {
	configRepo := repository.NewConfigRepo(cfg)
	socketRepo := repository.NewSocketRepo(cfg)
	return &handler{
		configRepo: configRepo,
		socketRepo: socketRepo,
	}
}
