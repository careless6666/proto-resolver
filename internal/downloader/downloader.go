package downloader

import (
	"ProtoDepsResolver/internal/models"
	"ProtoDepsResolver/internal/utils"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
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

	protoStorePath, err := utils.GetProtoStorePath()
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
				err = utils.CopyFileOrFolder(dep)
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

func GetRepoName(URL string) (string, error) {
	arr := strings.Split(URL, "/")
	lastPart := arr[len(arr)-1]
	if strings.HasSuffix(lastPart, ".git") {
		return lastPart[:len(lastPart)-4], nil
	}

	return "", errors.New("Invalid repo name: " + URL + " expected end with *.git")
}

func DownloadGitRepo(dep models.Dependency) error {

	gitInstalled := exec.Command("git", "--version")

	err := gitInstalled.Run()
	if err != nil {
		return errors.New("Git not installed")
	}

	protoStorePath, err := utils.GetAbsolutePathForDepInStore(dep)

	// TODO: problem repository with same name

	log.Println("git clone " + dep.Path + " to " + protoStorePath)

	if _, err := os.Stat(protoStorePath); !os.IsNotExist(err) {
		fmt.Println("repo exist skip clone")
		// TODO: git pull if flag enable update!
		//pullCmd := exec.Command("git", "clone", dep.Path, protoStorePath)
		//setStdCommand(pullCmd)
		//err = pullCmd.Run()

		//if err != nil {
		//	return err
		//}
	} else {
		cloneCmd := exec.Command("git", "clone", dep.Path, protoStorePath)

		setStdCommand(cloneCmd)
		err = cloneCmd.Run()

		if err != nil {
			return err
		}
	}

	gitFolder := path.Join(protoStorePath, ".git")

	if dep.Version.CommitRevision != "" {

		checkoutCmd := exec.Command("git", "--git-dir", gitFolder, "--work-tree", protoStorePath,
			"checkout", dep.Version.CommitRevision)

		setStdCommand(checkoutCmd)
		err = checkoutCmd.Run()

		if err != nil {
			return err
		}

	} else {
		cmd := exec.Command("git", "--git-dir", gitFolder, "--work-tree", protoStorePath,
			"checkout", "tags/"+dep.Version.Tag, "-f")

		setStdCommand(cmd)
		err = cmd.Run()

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func setStdCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}

func DownloadFile(dep models.Dependency) error {
	protoStorePath, err := utils.GetProtoStorePath()
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
