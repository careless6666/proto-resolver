package models

const (
	DependencyTypeGit  = iota
	DependencyTypeURL  = iota
	DependencyTypePath = iota
)

const RootPath string = "~/.proto-deps"

type Dependency struct {
	//locale | remote
	Type int
	Path string
	// for url | path
	DestinationPath string
	Version         *VersionInfo
}

type VersionInfo struct {
	Tag            string
	CommitRevision string
}
