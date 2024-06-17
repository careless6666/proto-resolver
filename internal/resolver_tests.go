package main

type IResolver interface {
	Resolve(dependency []Dependency) error
}

type Resolver struct {
}

func (Resolver) Resolve(dependency []Dependency) error {
	return nil
}
