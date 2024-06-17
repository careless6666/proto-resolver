package main

type IDepsFileParser interface {
	GetDeps(filePath string) ([]Dependency, error)
}

type Dependency struct {
	//locale | remote
	Type    string
	Path    string
	Version VersionInfo
}

type VersionInfo struct {
	Tag    string
	Commit string
}

type DepsFileParser struct {
}

func NewFileParser() IDepsFileParser {
	return DepsFileParser{}
}

func (f DepsFileParser) GetDeps(filePath string) ([]Dependency, error) {
	return nil, nil
}
