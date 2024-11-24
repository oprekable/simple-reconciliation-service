package root

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func Runner(cmd *cobra.Command, _ []string) (er error) {
	return cmd.Help()
}

func PersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	fromDate, errFrom := time.Parse("2006-01-02", FlagFromDateValue)

	if errFrom != nil {
		return fmt.Errorf("failed to parse flag 'from date': %v", FlagFromDateValue)
	}

	toDate, errTo := time.Parse("2006-01-02", FlagToDateValue)
	if errTo != nil {
		return fmt.Errorf("failed to parse flag 'to date': %v", FlagToDateValue)
	}

	if fromDate.After(toDate) {
		return fmt.Errorf("'from date': %v should before 'to date': %v", FlagFromDateValue, FlagToDateValue)
	}

	if FlagTotalDataSampleToGenerateValue < 0 {
		return fmt.Errorf("'amountdata': %v should bigger than 0", FlagTotalDataSampleToGenerateValue)
	}

	if FlagPercentageMatchSampleToGenerateValue < 0 || FlagPercentageMatchSampleToGenerateValue > 100 {
		return fmt.Errorf("'percentagematch': %v should between 0 and 100", FlagPercentageMatchSampleToGenerateValue)
	}

	return nil
}
