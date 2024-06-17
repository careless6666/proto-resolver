package main

import (
	"ProtoDepsResolver/cmd/app"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	util := cli.App{
		Name:  "protodep",
		Usage: "vendoring proto files with dependencies",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "gitlab_token",
				Value:    "",
				Usage:    "gitlab access token",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "github_token",
				Value:    "",
				Usage:    "github access token",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Action: app.Generate,
				Usage:  "generate proto files",
			},
		},
	}

	if err := util.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
