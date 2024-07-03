package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/careless6666/proto-resolver/internal/models"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var cachedStorePath *string = nil

func GetProtoStorePath() (string, error) {
	if cachedStorePath != nil {
		return *cachedStorePath, nil
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	result := path.Join(dirname, models.RootPath)
	cachedStorePath = &result

	return *cachedStorePath, nil
}

func GetRepoPathFromAddress(URL string) (string, error) {

	srcURL := URL

	if strings.HasPrefix(URL, "http://") {
		URL = URL[7:]
	}

	if strings.HasPrefix(URL, "https://") {
		URL = URL[8:]
	}

	if strings.HasSuffix(URL, ".git") {
		URL = URL[:len(URL)-4]
	}

	pathArr := strings.Split(URL, "/")

	if len(pathArr) < 2 {
		return "", errors.New("Invalid git repo path: " + srcURL)
	}

	return URL, nil
}

func GetRepoPathFromAddressInStorage(URL string) (string, error) {
	relativeAddress, err := GetRepoPathFromAddress(URL)

	if err != nil {
		return "", err
	}

	return relativeAddress, nil
}

func GetAbsolutePathForDepInStore(dep models.DependencyItem) (string, error) {
	protoStorePath, err := GetProtoStorePath()
	if err != nil {
		return "", err
	}

	switch dep.Type {
	case models.DependencyTypeGit:
		{
			relativePath, err := GetRepoPathFromAddressInStorage(dep.Source)
			if err != nil {
				return "", err
			}
			return path.Join(protoStorePath, relativePath), nil
		}

	case models.DependencyTypeURL:
		{
			projectId, err := GetProjectId()
			if err != nil {
				return "", err
			}
			dstFolderPath := path.Join(protoStorePath, projectId, dep.Type, dep.Tag, dep.RelativePath)
			return dstFolderPath, nil
		}
	case models.DependencyTypePath:
		{
			projectId, err := GetProjectId()
			if err != nil {
				return "", err
			}
			dstFolderPath := path.Join(protoStorePath, projectId, dep.Type, dep.Tag, dep.RelativePath)
			return dstFolderPath, nil
		}
	}

	return "", nil
}

func GetRelativePathForDepInStore(dep models.DependencyItem) (string, error) {
	switch dep.Type {
	case models.DependencyTypeGit:
		return GetRepoPathFromAddressInStorage(dep.Source)
	case models.DependencyTypeURL:
		{
			projectId, err := GetProjectId()
			if err != nil {
				return "", err
			}
			dstFolderPath := path.Join(projectId, dep.Type, dep.Tag, dep.RelativePath)
			return dstFolderPath, nil
		}
	case models.DependencyTypePath:
		{
			projectId, err := GetProjectId()
			if err != nil {
				return "", err
			}
			dstFolderPath := path.Join(projectId, dep.Type, dep.Tag, dep.RelativePath)
			return dstFolderPath, nil
		}
	}

	return "", nil
}

// To eliminate scenarios of colliding searches for nested proto files between projects with similar dependencies,
// let’s separate them by project separately
func GetProjectId() (string, error) {
	s, err := os.Getwd()
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(s))
	project_id := hex.EncodeToString(h.Sum(nil))
	return project_id, nil
}

func CopyFile(src, dst string) (err error) {
	dstDirectoryPath := filepath.Dir(dst)
	_, err = os.Stat(src)
	if err != nil {
		return
	}
	err = os.MkdirAll(dstDirectoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	defer out.Close()
	if err != nil {
		return err
	}
	input, err := os.OpenFile(src, 0, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, input)
	if err != nil {
		return err
	}

	return err
}

func CopyFileOrFolder(dep models.DependencyItem) error {
	//file or directory
	protoStorePath, err := GetProtoStorePath()
	if strings.HasSuffix(dep.RelativePath, ".proto") {
		file := filepath.Base(dep.RelativePath)

		if err != nil {
			return err
		}

		projectId, err := GetProjectId()
		if err != nil {
			return err
		}

		fullDstPath := filepath.Join(protoStorePath, projectId, dep.Type, dep.Tag, dep.RelativePath)

		err = os.MkdirAll(fullDstPath, os.ModePerm)
		if err != nil {
			return err
		}

		err = CopyFile(dep.RelativePath, path.Join(fullDstPath, file))
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

func CopyFilesRecursively(dep models.DependencyItem) error {
	protoStorePath, err := GetProtoStorePath()
	if err != nil {
		return err
	}
	return visitor("", dep, protoStorePath)
}

func visitor(currRelativePath string, dep models.DependencyItem, protoStorePath string) error {
	//copy files
	entries, err := os.ReadDir(path.Join(dep.Source, currRelativePath))
	if err != nil {
		return err
	}
	projectId, err := GetProjectId()
	if err != nil {
		return err
	}

	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".proto") && !e.Type().IsDir() {
			src := path.Join(dep.Source, currRelativePath, e.Name())
			dst := path.Join(protoStorePath, projectId, dep.Type, dep.Tag, dep.RelativePath, currRelativePath, e.Name())
			err = CopyFile(src, dst)
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
