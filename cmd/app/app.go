package app

import (
	"ProtoDepsResolver/internal/downloader"
	"ProtoDepsResolver/internal/parser"
	"ProtoDepsResolver/internal/resolver"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func New() (*App, error) {

	/*
			1) parse file
			2) download deps
		      a) copy from folder to deps folder for file deps
		      b) download file to deps folder from URL
		      c) clone git repo or update
			3) copy proto files to vendor.deps directory from home directory ~/.proto_deps
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

	err = resolver.Resolver{}.Resolve(deps)

	if err != nil {
		return err
	}

	return nil
}
