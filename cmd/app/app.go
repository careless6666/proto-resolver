package app

import (
	"ProtoDepsResolver/internal/downloader"
	"ProtoDepsResolver/internal/parser"
	"ProtoDepsResolver/internal/resolver"
	"ProtoDepsResolver/internal/utils"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strings"
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

	return &App{}, nil
}

type App struct{}

func Restore(ctx *cli.Context) error {

	if strings.EqualFold(ctx.String("verbose"), "true") {
		utils.Verbosity = true
	}

	fmt.Println("restored")

	var fileReader parser.IFileReader = parser.NewFileReader()
	depsParser := parser.NewFileParser(fileReader)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deps, err := depsParser.GetDeps(path.Join(pwd, "proto_deps.json"))

	if err != nil {
		return err
	}

	depsDownloader := downloader.NewDownloader(true)
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
