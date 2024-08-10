package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	util := cli.App{
		Name:  "proto-resolver",
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
			&cli.BoolFlag{
				Name:     "git_pull",
				Value:    true,
				Usage:    "git enable pull",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Value:    false,
				Usage:    "verbose",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "restore",
				Action: Restore,
				Usage:  "download all proto files",
			},
		},
	}

	if err := util.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
