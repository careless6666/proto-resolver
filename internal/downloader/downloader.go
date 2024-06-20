package downloader

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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
				err = DownloadFile(dep)
				if err != nil {
					return err
				}
				break
			}
		case models.DependencyTypeGit:
			{
				err = DownloadGitRepo(dep)
				if err != nil {
					return err
				}
				break
			}

		default:
			return errors.New("unknown dependency type, " + strconv.Itoa(dep.Type))
		}
	}

	return nil
}

func GetSshPathFromHttp(URL string) (string, error) {

	pathArr := strings.Split(URL, "/")

	if len(pathArr) < 3 {
		return "", errors.New("Invalid github repo path: " + URL)
	}

	sshPath := "git@" + URL + "/" + strings.Join(pathArr[1:], "/")

	return sshPath, nil
}

func DownloadGitRepo(dep models.Dependency) error {
	gitInstalled := exec.Command("git", "--version")

	protoStorePath, err := GetProtoStorePath()

	if err != nil {
		return err
	}

	err = gitInstalled.Run()
	if err != nil {
		return errors.New("Git not installed")
	}

	if dep.Version.CommitRevision != "" {
		
		address, err := GetSshPathFromHttp(dep.Path)
		if err != nil {
			return err
		}

		protoStorePath = path.Join(dep.Version.Tag, dep.Version.CommitRevision, protoStorePath)

		cloneCmd := exec.Command("git", "clone", address, protoStorePath, "--depth 1")

		fmt.Println(cloneCmd)

		return nil
	}

	return nil
}

func DownloadFile(dep models.Dependency) error {
	protoStorePath, err := GetProtoStorePath()
	urlArr := strings.Split(dep.Path, "/")
	dstFileName := urlArr[len(urlArr)-1]
	dstFilePath := path.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath)

	destinationFile := path.Join(dstFilePath, dstFileName)

	if _, err = os.Stat(destinationFile); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dstFilePath, os.ModePerm)
		if err != nil {
			return err
		}

		out, err := os.Create(destinationFile)
		defer out.Close()
		if err != nil {
			return err
		}
		resp, err := http.Get(dep.Path)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	return err
}

func copyFileOrFolder(dep models.Dependency) error {
	//file or directory
	protoStorePath, err := GetProtoStorePath()
	if strings.HasSuffix(dep.Path, ".proto") {
		file := filepath.Base(dep.Path)

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

	} else if strings.HasSuffix(dep.Path, "*") { //expected directory with one or many proto files
		dep.Path = dep.Path[:(len(dep.Path) - 1)]
		err := CopyFilesRecursively(dep)
		if err != nil {
			return err
		}
	} else { //expected directory with one or many proto files
		err := CopyFilesRecursively(dep)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyFilesRecursively(dep models.Dependency) error {
	protoStorePath, err := GetProtoStorePath()
	if err != nil {
		return err
	}
	return visitor("", dep, protoStorePath)
}

func visitor(currRelativePath string, dep models.Dependency, protoStorePath string) error {
	//copy files
	entries, err := os.ReadDir(path.Join(dep.Path, currRelativePath))
	if err != nil {
		return err
	}

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".proto") && !e.Type().IsDir() {
			src := path.Join(dep.Path, currRelativePath, e.Name())
			dst := path.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath, currRelativePath, e.Name())
			err = Copy(src, dst)
			if err != nil {
				return err
			}
		}

		//visit folders
		if e.Type().IsDir() {
			err = visitor(path.Join(currRelativePath, e.Name()), dep, protoStorePath)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func Copy(sourceFile, destinationFile string) (err error) {
	if _, err := os.Stat(destinationFile); errors.Is(err, os.ErrNotExist) {
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = ioutil.WriteFile(destinationFile, input, 0644)
		if err != nil {
			fmt.Println("Error creating", destinationFile)
			fmt.Println(err)
			return err
		}
		return err
	}

	return nil
}
