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
	Resolve(dependency []models.Dependency) error
}

type Resolver struct {
}

var _verbosity = false

func (Resolver) Resolve(dependency []models.Dependency) error {
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

func CopyProtoTree(dep models.Dependency) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	protoPath := path.Join(projectPath, vendorDeps, dep.DestinationPath)

	_, err = os.Stat(protoPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(protoPath, os.ModePerm); err != nil {
			return err
		}
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
		if !funk.Contains(file, dep.DestinationPath) {
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

func makeNewPathOnCopy(matchedFile string, dep models.Dependency) (string, error) {
	dstPath := dep.DestinationPath

	if dep.Type == models.DependencyTypeGit {
		pathFromAddress, err := utils.GetRepoPathFromAddress(dep.Path)
		if err != nil {
			return "", err
		}
		dstPath = path.Join(pathFromAddress, dep.DestinationPath)
	}

	start := funk.IndexOf(matchedFile, dstPath)

	//equal from middle to end
	if start+len(dstPath) == len(matchedFile) {
		return dstPath, nil
	}

	return dstPath + matchedFile[start+len(dstPath):], nil
}
