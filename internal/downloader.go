package main

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type IDownloader interface {
	Download(deps []models.Dependency) error
}

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(deps []models.Dependency) error {

	err := os.MkdirAll(models.RootPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		switch dep.Type {
		case models.DependencyTypePath:
			{
				err = copyFileOrFolder(dep, err)
				if err != nil {
					return err
				}
				break
			}
		case models.DependencyTypeURL:
			{
				break
			}
		case models.DependencyTypeGit:
			{
				break
			}

		default:
			return errors.New("unknown dependency type, " + strconv.Itoa(dep.Type))
		}
	}

	return nil

}

func copyFileOrFolder(dep models.Dependency, err error) error {
	//file or directory
	if strings.HasSuffix(dep.Path, ".proto") {
		err = os.MkdirAll(dep.DestinationPath, os.ModePerm)
		if err != nil {
			return err
		}
		r, err := os.Open(dep.Path)
		if err != nil {
			return err
		}
		defer r.Close()
		file := filepath.Base(dep.Path)
		w, err := os.Create(file)
		if err != nil {
			return err
		}
		defer w.Close()
		w.ReadFrom(r)
	} else if strings.HasSuffix(dep.Path, "*") {

	} else {

	}
	return nil
}
