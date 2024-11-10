package repository

import "github.com/google/wire"

func NewRepositories() *Repositories {
	return &Repositories{}
}

var Set = wire.NewSet(
	NewRepositories,
)
