package downloader

import (
	"ProtoDepsResolver/internal/models"
	"ProtoDepsResolver/internal/utils"
	"errors"
	"github.com/thoas/go-funk"
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

type Downloader struct {
	enablePull bool
}

func NewDownloader(enablePull bool) *Downloader {
	return &Downloader{
		enablePull: enablePull,
	}
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
				err = d.DownloadGitRepo(dep)
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

func (d *Downloader) DownloadGitRepo(dep models.Dependency) error {

	gitInstalled := exec.Command("git", "--version")

	err := gitInstalled.Run()
	if err != nil {
		return errors.New("Git not installed")
	}

	protoStorePath, err := utils.GetAbsolutePathForDepInStore(dep)

	// TODO: branch name as deps
	// TODO: git clone with github token
	// TODO: problem repository with same name
	// TODO: enable pull mode

	utils.LogVerbose("git clone " + dep.Path + " to " + protoStorePath)

	if _, err := os.Stat(protoStorePath); !os.IsNotExist(err) {
		utils.LogVerbose("repo exist skip clone")

		branchName, err := getBranchName(protoStorePath)
		if err != nil {
			return err
		}

		gitCheckoutCmd := exec.Command("git", "checkout", *branchName)
		gitCheckoutCmd.Dir = protoStorePath
		setStdCommand(gitCheckoutCmd)

		err = gitCheckoutCmd.Run()
		if err != nil {
			return err
		}

		if d.enablePull {
			pullCmd := exec.Command("git", "pull")
			pullCmd.Dir = path.Join(protoStorePath)
			setStdCommand(pullCmd)
			err = pullCmd.Run()

			if err != nil {
				return err
			}
		}

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

func getBranchName(protoStorePath string) (*string, error) {
	branchList := exec.Command("git", "branch", "-a")
	branchList.Dir = path.Join(protoStorePath)
	//setStdCommand(branchList)
	outputBranches, err := branchList.Output()
	if err != nil {
		return nil, err
	}

	utils.LogVerbose(string(outputBranches))
	branches := strings.Split(string(outputBranches), "\n")
	utils.LogVerbose(branches[0])
	for i := 0; i < len(branches); i++ {
		branches[i] = strings.Replace(branches[i], "*", "", -1)
		branches[i] = strings.TrimSpace(branches[i])
	}

	//search with priority
	// 1) master or main
	// 2) any local
	// 3) any remote

	for _, branch := range branches {
		if (funk.Contains(branch, "main") || funk.Contains(branch, "master")) && !funk.Contains(branch, "remotes") {
			return &branch, nil
		}
	}

	for _, branch := range branches {
		if !funk.Contains(branch, "remotes") && !funk.Contains(branch, "detached") {
			return &branch, nil
		}
	}

	for _, branch := range branches {
		if funk.Contains(branch, "remotes") {
			return &branch, nil
		}
	}

	return nil, errors.New("Could not find any branch to fetch")
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
