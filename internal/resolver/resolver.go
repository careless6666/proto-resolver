package resolver

import (
	"ProtoDepsResolver/internal/models"
	"ProtoDepsResolver/internal/utils"
	"github.com/mattn/go-zglob"
	"github.com/thoas/go-funk"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
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

	repoPath, err := utils.GetRelativePathForDepInStore(dep)
	if err != nil {
		return err
	}

	protoPath := path.Join(projectPath, vendorDeps, repoPath)

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

	fullDstPath := ""
	if dep.Type == models.DependencyTypePath || dep.Type == models.DependencyTypeURL {
		fullDstPath = filepath.Join(protoStorePath, dep.Version.Tag, dep.DestinationPath)
	} else {
		fullDstPath, err = utils.GetRepoPathFromAddress(dep.Path)
		if err != nil {
			return err
		}

		fullDstPath = path.Join(fullDstPath, dep.GitPath)
	}

	for _, file := range matches {
		logInfo(file)
		if !funk.Contains(file, fullDstPath) {
			continue
		}
		err = CopyFile(file, path.Join(projectPath, vendorDeps, fullDstPath))
		if err != nil {
			return err
		}
	}

	return nil
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

func logInfo(message string) {
	if _verbosity {
		log.Println(message)
	}
}
