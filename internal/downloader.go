package main

type IDownloader interface {
	Download(deps []Dependency) error
}

type Downloader struct{}

func (d Downloader) Download(deps []Dependency) error {
	return nil
}
