package handler

import (
	"context"
	"fmt"

	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/repository"
)

type Handler interface {
	Handle(ctx context.Context, data string) (string, error)
}

type handler struct {
	socketRepo repository.SocketRepo
}

func (h *handler) Handle(ctx context.Context, data string) (string, error) {
	if data == "port-start" {
		localPort := 3307
		remoteHost := "localhost"
		remotePort := 3306

		host := models.SSHConfig{
			Name:     "FirstHost",
			Host:     "<public-host-ip>",
			Port:     22,
			User:     "user1",
			Password: "samplepassword",
			NextHop: &models.SSHConfig{
				Name:    "FinalHost",
				Host:    "<final-private-host-ip>",
				Port:    22,
				User:    "user2",
				KeyFile: "~/.ssh/.id_rsa",
			},
		}
		go h.socketRepo.StartConnection(host, localPort, remoteHost, remotePort)
		return fmt.Sprintf("Forwarding local port %d to %s:%d via SSH", localPort, remoteHost, remotePort), nil
	}
	if data == "port-stop" {
		localPort := 3307
		remoteHost := "localhost"
		remotePort := 3306
		h.socketRepo.StopConnection()
		return fmt.Sprintf("Forwarding local port %d to %s:%d via SSH", localPort, remoteHost, remotePort), nil
	}
	return "Hello" + data, nil
}

func NewHandler() Handler {
	socketRepo := repository.NewSocketRepo()

	return &handler{
		socketRepo: socketRepo,
	}
}
