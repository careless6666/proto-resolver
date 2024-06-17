package app

import (
	"github.com/urfave/cli/v2"
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

func Generate(ctx *cli.Context) error {
	/*a, err := New(ctx.String("gitlab_token"), ctx.String("github_token"))
	if err != nil {
		log.Fatal(err)
	}*/
	//dependencies, err := a.protodepProvider.FileDependencies()
	//if err != nil {
	//	log.Fatal(err)
	//	}
	//if len(dependencies) == 0 {
	//	return errors.New("empty dependencies")
	//}

	// return a.dependencyManager.Process(ctx.Context, dependencies...)
	return nil
}
