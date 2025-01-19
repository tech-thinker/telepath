package main

import (
	"fmt"
	"os"

	"github.com/tech-thinker/telepath/cmd"
	"github.com/tech-thinker/telepath/config"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/daemon"
	"github.com/tech-thinker/telepath/handler"
	"github.com/urfave/cli/v2"
)

var (
	AppVersion = "v0.0.0"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

func main() {
	cfg := config.InitConfig()
	handler := handler.NewHandler(cfg)
	daemonMgr := daemon.NewDaemonMgr(constants.PID_FILE_PATH, constants.SOCKET_PATH, handler)
	appCmd := cmd.NewApp(daemonMgr)

	app := &cli.App{
		Name:    "telepath",
		Version: AppVersion,
		Commands: []*cli.Command{
			appCmd.Daemon(),
			appCmd.Crediential(),
			appCmd.Host(),
			appCmd.Tunnel(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
