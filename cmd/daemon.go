package cmd

import (
	"context"

	"github.com/urfave/cli/v2"
)

func (a *app) Daemon() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "crediential operation",
		Subcommands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start daemon",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:   "daemon-child",
						Usage:  "Run the service as child daemon",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					if c.Bool("daemon-child") {
						return a.daemonMgr.RunDaemonChild(ctx)
					}
					return a.daemonMgr.RunAsDaemon(ctx)
				},
			},
			{
				Name:  "stop",
				Usage: "Stop daemon",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					return a.daemonMgr.StopDaemon(ctx)
				},
			},
			{
				Name:  "status",
				Usage: "Status of daemon",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					return a.daemonMgr.StatusDaemon(ctx)
				},
			},
		},
	}
}
