package main

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"strconv"
)

type IDownloader interface {
	Download(deps []models.Dependency) error
}

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(deps []models.Dependency) error {

	for _, dep := range deps {
		switch dep.Type {
		case models.DependencyTypePath:
			{

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
