package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/urfave/cli/v2"
)

func (a *app) Tunnel() *cli.Command {
	return &cli.Command{
		Name:  "tunnel",
		Usage: "Tunnel operation",
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add Tunnel",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "local-port",
						Aliases:  []string{"L"},
						Usage:    "Local port.",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "remote-host",
						Aliases: []string{"H"},
						Usage:   "Remote host.",
					},
					&cli.IntFlag{
						Name:     "remote-port",
						Aliases:  []string{"R"},
						Usage:    "Remote port.",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "host-chain",
						Aliases: []string{"C"},
						Usage:   "List of host name as jump host.",
					},
				},
				Action: func(c *cli.Context) error {
					localPort := c.Int("local-port")
					remoteHost := c.String("remote-host")
					remotePort := c.Int("remote-port")
					hosts := c.String("host-chain")

					tunnel := models.Tunnel{
						Name:       c.Args().First(),
						LocalPort:  localPort,
						RemoteHost: remoteHost,
						RemotePort: remotePort,
						HostChain:  strings.Split(hosts, ","),
					}

					packet := models.Packet{
						Action: "add-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
						Data:   tunnel.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove Tunnel",
				Action: func(c *cli.Context) error {
					tunnel := models.Tunnel{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "remove-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
						Data:   tunnel.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "detail",
				Usage: "Detail Tunnel",
				Action: func(c *cli.Context) error {
					tunnel := models.Tunnel{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "detail-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
						Data:   tunnel.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List Tunnel",
				Action: func(c *cli.Context) error {
					packet := models.Packet{
						Action: "list-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "start",
				Usage: "Start Tunnel",
				Action: func(c *cli.Context) error {
					tunnel := models.Tunnel{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "start-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
						Data:   tunnel.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "Stop Tunnel",
				Action: func(c *cli.Context) error {
					tunnel := models.Tunnel{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "stop-tunnel",
						Type:   constants.PACKET_TYPE_TUNNEL,
						Data:   tunnel.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(context.Background(), packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
		},
	}
}
