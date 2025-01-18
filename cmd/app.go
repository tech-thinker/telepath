package cmd

import (
	"context"

	"github.com/tech-thinker/telepath/daemon"
	"github.com/urfave/cli/v2"
)

type App interface {
	StartDaemon() *cli.Command
	StopDaemon() *cli.Command
	SendCommand() cli.ActionFunc
}

type app struct {
	daemonMgr daemon.DaemonMgr
}

func (a *app) StartDaemon() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start the daemon service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run the service as a daemon",
			},
			&cli.BoolFlag{
				Name:   "daemon-child",
				Usage:  "Run the service as child daemon",
				Hidden: true,
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			if c.Bool("daemon") {
				return a.daemonMgr.RunAsDaemon(ctx)
			}
			if c.Bool("daemon-child") {
				return a.daemonMgr.RunDaemonChild(ctx)
			}
			return a.daemonMgr.SendCommandToDaemon(ctx, c.Args().Slice())
		},
	}
}

func (a *app) StopDaemon() *cli.Command {
	return &cli.Command{
		Name:  "kill",
		Usage: "Start the daemon service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run the service as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			if c.Bool("daemon") {
				return a.daemonMgr.StopDaemon(ctx)
			}
			return a.daemonMgr.SendCommandToDaemon(ctx, c.Args().Slice())
		},
	}
}

func (a *app) SendCommand() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return a.daemonMgr.SendCommandToDaemon(context.Background(), ctx.Args().Slice())
	}
}

func NewApp(
	daemonMgr daemon.DaemonMgr,
) App {
	return &app{
		daemonMgr: daemonMgr,
	}
}
