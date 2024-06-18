package main

import "ProtoDepsResolver/internal/parser"

type IResolver interface {
	Resolve(dependency []parser.Dependency) error
}

type Resolver struct {
}

func (Resolver) Resolve(dependency []parser.Dependency) error {
	return nil
}
