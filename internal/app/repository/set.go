package repository

import (
	"simple-reconciliation-service/internal/app/repository/sample"

	"github.com/google/wire"
)

func NewRepositories(
	repoSample sample.Repository,
) *Repositories {
	return &Repositories{
		RepoSample: repoSample,
	}
}

var Set = wire.NewSet(
	sample.Set,
	NewRepositories,
)
