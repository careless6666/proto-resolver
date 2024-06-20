package downloader

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func GetProtoStorePath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dirname, models.RootPath), nil

}

func (d *Downloader) Download(deps []models.Dependency) error {

	protoStorePath, err := GetProtoStorePath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(protoStorePath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		switch dep.Type {
		case models.DependencyTypePath:
			{
				err = copyFileOrFolder(dep)
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

func copyFileOrFolder(dep models.Dependency) error {
	//file or directory
	if strings.HasSuffix(dep.Path, ".proto") {
		file := filepath.Base(dep.Path)
		protoStorePath, err := GetProtoStorePath()
		if err != nil {
			return err
		}

		fullDstPath := filepath.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath)

		err = os.MkdirAll(fullDstPath, os.ModePerm)
		if err != nil {
			return err
		}

		err = Copy(dep.Path, path.Join(fullDstPath, file))
		if err != nil {
			return err
		}

	} else if strings.HasSuffix(dep.Path, "*") {

	} else {

	}
	return nil
}

func Copy(sourceFile, destinationFile string) (err error) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
	return err
}
