package utils

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"os"
	"path"
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
