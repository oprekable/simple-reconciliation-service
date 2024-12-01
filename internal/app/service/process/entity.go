package process

type ReconciliationSummary struct {
	FileMissingBankTrx              map[string]string
	FileMissingSystemTrx            string
	FileMatchedSystemTrx            string
	TotalProcessedSystemTrx         int64
	TotalMatchedSystemTrx           int64
	TotalNotMatchedSystemTrx        int64
	SumAmountProcessedSystemTrx     float64
	SumAmountMatchedSystemTrx       float64
	SumAmountDiscrepanciesSystemTrx float64
}
