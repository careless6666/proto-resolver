package models

const (
	DependencyTypeGit  = iota
	DependencyTypeURL  = iota
	DependencyTypePath = iota
)

const RootPath string = ".proto_deps"

type Dependency struct {
	//locale | remote
	Type int
	Path string
	// path where should be stored inside result repository
	DestinationPath string
	Version         *VersionInfo
}

type VersionInfo struct {
	Tag            string
	CommitRevision string
}

type DependencyList struct {
	Type           string `json:"type"`
	Source         string `json:"source"`
	RelativePath   string `json:"RelativePath"`
	Version        string `json:"version"`
	Tag            string `json:"tag"`
	Branch         string `json:"branch"`
	CommitRevision string `json:"commitRevision"`
}

type DependencyRoot struct {
	Version string         `json:"version"`
	Deps    DependencyList `json:"deps"`
}
