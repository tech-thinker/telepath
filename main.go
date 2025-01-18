package main

import (
	"fmt"
	"os"

	"github.com/tech-thinker/telepath/cmd"
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
	handler := handler.NewHandler()
	daemonMgr := daemon.NewDaemonMgr(constants.PID_FILE_PATH, constants.SOCKET_PATH, handler)
	appCmd := cmd.NewApp(daemonMgr)

	app := &cli.App{
		Name:    "protty",
		Version: AppVersion,
		Action:  appCmd.SendCommand(),
		Commands: []*cli.Command{
			appCmd.StartDaemon(),
			appCmd.StopDaemon(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
