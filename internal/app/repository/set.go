package repository

import (
	"simple-reconciliation-service/internal/app/repository/process"
	"simple-reconciliation-service/internal/app/repository/sample"

	"github.com/google/wire"
)

func NewRepositories(
	repoSample sample.Repository,
	repoProcess process.Repository,
) *Repositories {
	return &Repositories{
		RepoSample:  repoSample,
		RepoProcess: repoProcess,
	}
}

var Set = wire.NewSet(
	sample.Set,
	process.Set,
	NewRepositories,
)
