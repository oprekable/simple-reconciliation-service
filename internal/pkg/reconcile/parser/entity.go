package parser

import (
	"simple-reconciliation-service/internal/pkg/reconcile/parser/banks"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
)

type TrxData struct {
	SystemTrx       []*systems.SystemTrxData
	BankTrx         []*banks.BankTrxData
	MinSystemAmount float64
	MaxSystemAmount float64
}
