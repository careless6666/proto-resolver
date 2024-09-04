package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	apiToken := ""
	util := cli.App{
		Name:  "proto-resolver",
		Usage: "vendoring proto files with dependencies",
		Commands: []*cli.Command{
			{
				Name:   "restore",
				Action: Restore,
				Usage:  "download all proto files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "gitlab_token",
						Aliases:     []string{"a"},
						Usage:       "gitlab access token",
						Required:    false,
						Destination: &apiToken,
					},
					&cli.StringFlag{
						Name:     "github_token",
						Usage:    "github access token",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "gitlab_domain",
						Usage:    "gitlab special domain",
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
			},
		},
	}

	if err := util.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
