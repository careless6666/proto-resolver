package app

import (
	"ProtoDepsResolver/internal/downloader"
	"ProtoDepsResolver/internal/parser"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func New(gitlabToken, githubToken string) (*App, error) {

	/*
		1) parse file
		2) download deps
		3) copy proto files  vendor.deps
	*/
	//var err error
	/*
		var protodepProvider *protodep.Provider
		{
			protodepProvider = protodep.New()
		}

		var githubProvider *github.Provider
		{
			githubProvider, err = github.New(context.Background(), github.Config{
				Token: githubToken,
			})
			if err != nil {
				return nil, err
			}
		}

		var gitlabProvider *gitlab.Provider
		{
			gitlabProvider, err = gitlab.New(gitlab.Config{
				Token: gitlabToken,
			})
			if err != nil {
				return nil, err
			}
		}

		var aggregator *provider.Aggregator
		{
			aggregator = provider.New(map[model.Domain]model.Provider{
				model.DomainGitlab: gitlabProvider,
				model.DomainGithub: githubProvider,
			})
		}

		var relativeProvider *relative.Provider
		{
			relativeProvider = relative.New(aggregator, model.DomainGithub, model.DomainGitlab)
		}

		var wktProvider *wkt.Provider
		{
			wktProvider = wkt.New(relativeProvider)
		}

		var dependencyManager *dependency.Manager
		{
			dependencyManager = dependency.New(wktProvider)
		}
	*/

	return &App{
		//protodepProvider:  protodepProvider,
		//dependencyManager: dependencyManager,
	}, nil
}

type App struct {
	//protodepProvider  *protodep.Provider
	//dependencyManager *dependency.Manager
}

func Restore(ctx *cli.Context) error {
	fmt.Println("restored")

	var fileReader parser.IFileReader = parser.NewFileReader()
	depsParser := parser.NewFileParser(fileReader)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deps, err := depsParser.GetDeps(pwd + "/proto_deps.yml")

	if err != nil {
		return err
	}

	depsDownloader := downloader.NewDownloader()
	err = depsDownloader.Download(deps)

	if err != nil {
		return err
	}

	return nil
}
