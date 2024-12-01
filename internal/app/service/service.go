package service

import (
	"simple-reconciliation-service/internal/app/service/process"
	"simple-reconciliation-service/internal/app/service/sample"
)

type Services struct {
	SvcSample  sample.Service
	SvcProcess process.Service
}
