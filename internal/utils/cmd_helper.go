package utils

import (
	"github.com/careless6666/proto-resolver/internal/models"
	"github.com/urfave/cli/v2"
)

func ReadOptions(ctx *cli.Context) models.CmdOptions {

	opt := models.CmdOptions{}

	ctx.String("verbose")
	opt.Verbose = ctx.Bool("verbose")
	opt.GitPull = ctx.Bool("git-pull")
	opt.GitlabToken = ctx.String("gitlab_token")
	opt.GithubToken = ctx.String("github_token")
	opt.GitlabDomain = ctx.String("gitlab_domain")

	return opt
}
