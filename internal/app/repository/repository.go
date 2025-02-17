package repository

import (
	"simple-reconciliation-service/internal/app/repository/process"
	"simple-reconciliation-service/internal/app/repository/sample"
)

type Repositories struct {
	RepoSample  sample.Repository
	RepoProcess process.Repository
}
