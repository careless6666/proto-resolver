package models

const (
	DependencyTypeGit  = "git"
	DependencyTypeURL  = "url"
	DependencyTypePath = "path"
)

const RootPath string = ".proto_deps"

type DependencyItem struct {
	Type         string `json:"type"`
	Source       string `json:"source"`
	RelativePath string `json:"relativePath"`
	//for url anf file types only
	Version        string `json:"version"`
	Tag            string `json:"tag"`
	Branch         string `json:"branch"`
	CommitRevision string `json:"commitRevision"`
	//	ExcludeRegex   string `json:"exclude_regex"`
}

type DependencyRoot struct {
	Version string           `json:"version"`
	Deps    []DependencyItem `json:"deps"`
}

type CmdOptions struct {
	Verbose      bool
	GitPull      bool
	GithubToken  string
	GitlabToken  string
	GitlabDomain string
}
