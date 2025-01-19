package cmd

import (
	"context"
	"fmt"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/urfave/cli/v2"
)

func (a *app) Host() *cli.Command {
	return &cli.Command{
		Name:  "host",
		Usage: "Host operation",
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add Host",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "host",
						Aliases:  []string{"H"},
						Usage:    "SSH server address.",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"P"},
						Usage:   "SSH server port.",
					},
					&cli.StringFlag{
						Name:    "user",
						Aliases: []string{"U"},
						Usage:   "SSH server user name.",
					},
					&cli.StringFlag{
						Name:    "cred",
						Aliases: []string{"C"},
						Usage:   "Crediential name.",
					},
				},
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					host := c.String("host")
					port := c.Int("port")
					user := c.String("user")
					cred := c.String("cred")

					hostCfg := models.HostConfig{
						Name:            c.Args().First(),
						Host:            host,
						Port:            port,
						User:            user,
						CredientialName: cred,
					}

					packet := models.Packet{
						Action: "add-host",
						Type:   constants.PACKET_TYPE_HOST_CONFIG,
						Data:   hostCfg.ToByte(),
					}
					err := a.daemonMgr.SendCommandToDaemon(ctx, packet)
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove Host",
				Action: func(c *cli.Context) error {
					hostCfg := models.HostConfig{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "remove-host",
						Type:   constants.PACKET_TYPE_HOST_CONFIG,
						Data:   hostCfg.ToByte(),
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
				Usage: "Detail Host",
				Action: func(c *cli.Context) error {
					hostCfg := models.HostConfig{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "detail-host",
						Type:   constants.PACKET_TYPE_HOST_CONFIG,
						Data:   hostCfg.ToByte(),
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
				Usage: "List Host",
				Action: func(c *cli.Context) error {
					packet := models.Packet{
						Action: "list-host",
						Type:   constants.PACKET_TYPE_HOST_CONFIG,
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
