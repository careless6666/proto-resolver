package main

import "ProtoDepsResolver/internal/parser"

type IDownloader interface {
	Download(deps []parser.Dependency) error
}

type Downloader struct{}

func (d Downloader) Download(deps []parser.Dependency) error {
	return nil
}
