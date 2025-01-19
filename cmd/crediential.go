package cmd

import (
	"context"
	"fmt"

	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/urfave/cli/v2"
)

func (a *app) Crediential() *cli.Command {
	return &cli.Command{
		Name:  "crediential",
		Usage: "crediential operation",
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add crediential",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "type",
						Aliases:  []string{"T"},
						Usage:    "Crediential type [PASS/KEY].",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "pass",
						Aliases: []string{"P"},
						Usage:   "Plain text password.",
					},
					&cli.StringFlag{
						Name:    "key-file",
						Aliases: []string{"K"},
						Usage:   "SSH Key file path.",
					},
				},
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					credType := c.String("type")
					pass := c.String("pass")
					key := c.String("key-file")

					cred := models.Crediential{
						Name:     c.Args().First(),
						Type:     credType,
						Password: pass,
						KeyFile:  key,
					}

					packet := models.Packet{
						Action: "add-cred",
						Type:   constants.PACKET_TYPE_CREDIENTIAL,
						Data:   cred.ToByte(),
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
				Usage: "Remove crediential",
				Action: func(c *cli.Context) error {
					cred := models.Crediential{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "remove-cred",
						Type:   constants.PACKET_TYPE_CREDIENTIAL,
						Data:   cred.ToByte(),
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
				Usage: "Detail crediential",
				Action: func(c *cli.Context) error {
					cred := models.Crediential{
						Name: c.Args().First(),
					}

					packet := models.Packet{
						Action: "detail-cred",
						Type:   constants.PACKET_TYPE_CREDIENTIAL,
						Data:   cred.ToByte(),
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
				Usage: "List crediential",
				Action: func(c *cli.Context) error {
					packet := models.Packet{
						Action: "list-cred",
						Type:   constants.PACKET_TYPE_CREDIENTIAL,
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
