package utils

import (
	"ProtoDepsResolver/internal/models"
	"errors"
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

func GetAbsolutePathForDepInStore(dep models.Dependency) (string, error) {
	protoStorePath, err := GetProtoStorePath()
	if err != nil {
		return "", err
	}

	switch dep.Type {
	case models.DependencyTypeGit:
		{
			relativePath, err := GetRepoPathFromAddressInStorage(dep.Path)
			if err != nil {
				return "", err
			}
			return path.Join(protoStorePath, relativePath), nil
		}

	case models.DependencyTypeURL:
		{
			dstFolderPath := path.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath)
			return dstFolderPath, nil
		}
	case models.DependencyTypePath:
		{
			dstFolderPath := path.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath)
			return dstFolderPath, nil
		}
	}

	return "", nil
}

func GetRelativePathForDepInStore(dep models.Dependency) (string, error) {
	switch dep.Type {
	case models.DependencyTypeGit:
		return GetRepoPathFromAddressInStorage(dep.Path)
	case models.DependencyTypeURL:
		{
			dstFolderPath := path.Join(dep.Version.Tag, dep.DestinationPath)
			return dstFolderPath, nil
		}
	case models.DependencyTypePath:
		{
			dstFolderPath := path.Join(dep.Version.Tag, dep.DestinationPath)
			return dstFolderPath, nil
		}
	}

	return "", nil
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

func CopyFileOrFolder(dep models.Dependency) error {
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

		err = CopyFile(dep.Path, path.Join(fullDstPath, file))
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
