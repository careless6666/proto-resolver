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
