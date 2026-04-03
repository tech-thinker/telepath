package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tech-thinker/telepath/services"
	"github.com/urfave/cli/v2"
)

var (
	AppVersion = "v0.0.0"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

func main() {
	var configPath string
	var dryRun bool = false

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version: %s\n", AppVersion)
		fmt.Printf("Commit: %s\n", CommitHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
	}

	app := &cli.App{
		Name:        "telepath",
		Version:     AppVersion,
		Description: "telepath is a cli application to forward port securly.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"f"},
				Usage:       "accepts config file path.",
				Required:    true,
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "use to tests your configuration.",
				Required:    false,
				Destination: &dryRun,
			},
		},
		Action: func(ctx *cli.Context) error {
			cfgs, err := services.ParseConfig(configPath)
			if err != nil {
				return err
			}
			cfgs, err = services.ValidateConfig(cfgs)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println("Your configuration format is valid.")
				return nil
			}

			// Root context — cancelled on Ctrl+C / SIGTERM
			rootCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Listen for shutdown signals
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			go func() {
				sig := <-sigCh
				log.Printf("Received signal %s, shutting down gracefully...", sig)
				cancel()
			}()

			s := services.NewServer(cfgs)
			var wg sync.WaitGroup
			s.StartAll(rootCtx, &wg)

			// Block until all tunnels have exited cleanly
			wg.Wait()
			log.Println("All tunnels stopped. Goodbye.")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
