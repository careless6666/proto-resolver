package resolver

import (
	"ProtoDepsResolver/internal/models"
	"ProtoDepsResolver/internal/utils"
	"github.com/mattn/go-zglob"
	"github.com/thoas/go-funk"
	"os"
	"path"
)

const (
	vendorDeps = "vendor.pb"
)

type IResolver interface {
	Resolve(dependency []models.DependencyItem) error
}

type Resolver struct {
}

var _verbosity = false

func (Resolver) Resolve(dependency []models.DependencyItem) error {
	if err := os.RemoveAll(path.Join(vendorDeps)); err != nil {
		return err
	}

	for _, dep := range dependency {
		err := CopyProtoTree(dep)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyProtoTree(dep models.DependencyItem) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	protoStorePath, err := utils.GetProtoStorePath()
	if err != nil {
		return err
	}

	restoreRelativePath, err := utils.GetRelativePathForDepInStore(dep)
	if err != nil {
		return err
	}

	absolutePath := path.Join(protoStorePath, restoreRelativePath)

	matches, err := zglob.Glob(absolutePath + "/**/*.proto")
	if err != nil {
		return err
	}

	_, err = os.Stat(path.Join(projectPath, vendorDeps))

	if os.IsNotExist(err) {
		if err = os.MkdirAll(projectPath+"/"+vendorDeps, os.ModePerm); err != nil {
			return err
		}
	}

	for _, file := range matches {
		utils.LogInfo(file)
		if !funk.Contains(file, dep.RelativePath) {
			continue
		}

		relativePath, err := makeNewPathOnCopy(file, dep)
		if err != nil {
			return err
		}

		fullDstPath := path.Join(projectPath, vendorDeps, relativePath)
		err = utils.CopyFile(file, fullDstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func makeNewPathOnCopy(matchedFile string, dep models.DependencyItem) (string, error) {
	dstPath := dep.RelativePath

	if dep.Type == models.DependencyTypeGit {
		pathFromAddress, err := utils.GetRepoPathFromAddress(dep.Source)
		if err != nil {
			return "", err
		}
		dstPath = path.Join(pathFromAddress, dep.RelativePath)
	}

	start := funk.IndexOf(matchedFile, dstPath)

	//equal from middle to end
	if start+len(dstPath) == len(matchedFile) {
		return dstPath, nil
	}

	return dstPath + matchedFile[start+len(dstPath):], nil
}
