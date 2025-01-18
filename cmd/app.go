package cmd

import (
	"github.com/tech-thinker/telepath/daemon"
	"github.com/urfave/cli/v2"
)

type App interface {
	Daemon() *cli.Command
	Crediential() *cli.Command
	Host() *cli.Command
	Tunnel() *cli.Command
}

type app struct {
	daemonMgr daemon.DaemonMgr
}

func NewApp(
	daemonMgr daemon.DaemonMgr,
) App {
	return &app{
		daemonMgr: daemonMgr,
	}
}
