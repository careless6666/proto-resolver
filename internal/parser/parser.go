//go:generate mockgen -source=parser.go -destination=mock/parser-mock.go -package mock

package parser

import (
	"encoding/json"
	"errors"
	"github.com/careless6666/proto-resolver/internal/models"
	"os"
	"regexp"
)

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

func (f *DepsFileParser) GetDeps(path string) ([]models.DependencyItem, error) {
	content, err := f.fileReader.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("no dependencies found, empty file")
	}

	targets := models.DependencyRoot{}

	err = json.Unmarshal(content, &targets)
	if err != nil {
		return nil, err
	}

	//check version
	matchedVersion, err := regexp.Match(`(v\d+)`, []byte(targets.Version))
	if err != nil {
		return nil, err
	}

	if !matchedVersion {
		return nil, errors.New("invalid dependencies found, not a version")
	}

	if len(targets.Deps) < 1 {
		return nil, errors.New("invalid dependencies file, no deps found")
	}

	return targets.Deps, nil
}
