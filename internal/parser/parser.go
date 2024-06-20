//go:generate mockgen -source=parser.go -destination=mock/parser-mock.go -package mock

package parser

import (
	"ProtoDepsResolver/internal/models"
	"errors"
	"os"
	"regexp"
	"strings"
)

//type IDepsFileParser interface {
//	GetDeps(path string) ([]Dependency, error)
//}

type IFileReader interface {
	ReadFile(filePath string) ([]byte, error)
}

type FileReader struct {
}

func (f *FileReader) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func NewFileReader() *FileReader {
	return &FileReader{}
}

type DepsFileParser struct {
	fileReader IFileReader
}

func NewFileParser(fileReader IFileReader) *DepsFileParser {
	return &DepsFileParser{
		fileReader: fileReader,
	}
}

func (f *DepsFileParser) GetDeps(path string) ([]models.Dependency, error) {
	content, err := f.fileReader.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("no dependencies found, empty file")
	}

	depsStr := strings.Split(string(content), "\n")
	//check version
	matchedVersion, err := regexp.Match(`(version)(:)( )(v\d+)`, []byte(depsStr[0]))
	if err != nil {
		return nil, err
	}

	if !matchedVersion {
		return nil, errors.New("invalid dependencies found, not a version")
	}

	if len(depsStr) < 3 {
		return nil, errors.New("invalid dependencies file, rows count less than 3")
	}
	//check deps
	if depsStr[1] != "deps:" {
		return nil, errors.New("invalid dependencies file, \"deps:\" block not found")
	}

	result := make([]models.Dependency, 0)
	//get deps
	for _, depStr := range depsStr[2:] {
		dep, err := ParseDepsLine(depStr)
		if err != nil {

			return nil, err
		}
		if dep == nil {
			continue
		}

		result = append(result, *dep)
	}

	if len(result) == 0 {
		return nil, errors.New("no dependencies found")
	}

	return result, nil
}

func ParseDepsLine(dependency string) (*models.Dependency, error) {
	dependency = strings.TrimSpace(dependency)

	// skip empty lines
	if dependency == "" {
		return nil, nil
	}

	// skip commented line
	if strings.HasPrefix(dependency, "#") {
		return nil, nil
	}

	matchedGit, err := regexp.Match(`- git: `, []byte(dependency))
	if err != nil {
		return nil, err
	}

	if matchedGit {
		return getGitDeps(dependency[7:])
	}

	matchedURL, err := regexp.Match(`- url: `, []byte(dependency))
	if err != nil {
		return nil, err
	}

	if matchedURL {
		return getUrlDeps(dependency[7:])
	}

	matchedFile, err := regexp.Match(`- path: `, []byte(dependency))
	if err != nil {
		return nil, err
	}

	if matchedFile {
		return getFileDeps(dependency[8:])
	}

	return nil, nil
}

func getGitDeps(dependency string) (*models.Dependency, error) {
	depPaths := strings.Split(dependency, " ")
	if len(depPaths) != 3 {
		return nil, errors.New("invalid dependency, have to by pattern \"- git: github.com/repo/file.proto v0.0.0-20211005231101-409e134ffaac\"")
	}

	version := &models.VersionInfo{
		Tag:            "",
		CommitRevision: "",
	}

	versionStr := strings.Split(depPaths[2], "-")
	if len(versionStr) == 1 {
		version.Tag = versionStr[0]
	} else if len(versionStr) == 3 {
		version.Tag = versionStr[0]
		version.CommitRevision = versionStr[2]
	} else {
		return nil, errors.New("invalid count items in git deps version info")
	}

	return &models.Dependency{
		Type:            models.DependencyTypeGit,
		Path:            depPaths[0],
		GitPath:         depPaths[1],
		DestinationPath: "",
		Version:         version,
	}, nil
}

func getUrlDeps(dependency string) (*models.Dependency, error) {
	depPaths := strings.Split(dependency, " ")
	if len(depPaths) != 3 {
		return nil, errors.New("invalid dependency, have to by pattern \"- url: https://github.com/repo/file.proto ./github.com/repo/file.proto v1\"")
	}

	matchedProtoFileURL, err := regexp.Match(`(http:\/\/)(.*)(.)(proto)|(https:\/\/)(.*)(.)(proto)`, []byte(depPaths[0]))
	if err != nil {
		return nil, err
	}

	if matchedProtoFileURL {
		return &models.Dependency{
			Type:            models.DependencyTypeURL,
			Path:            depPaths[0],
			DestinationPath: depPaths[1],
			Version: &models.VersionInfo{
				Tag:            depPaths[2],
				CommitRevision: "",
			},
		}, nil
	}

	return nil, errors.New("invalid dependency, expected URL to proto file")
}

func getFileDeps(dependency string) (*models.Dependency, error) {
	depPaths := strings.Split(dependency, " ")
	if len(depPaths) != 3 {
		return nil, errors.New("invalid dependency, have to by pattern \"- path: /var/github.com/repo/file.proto ./github.com/repo/file.proto v1\"")
	}

	return &models.Dependency{
		Type:            models.DependencyTypePath,
		Path:            depPaths[0],
		DestinationPath: depPaths[1],
		Version: &models.VersionInfo{
			Tag: depPaths[2],
		},
	}, nil
}
