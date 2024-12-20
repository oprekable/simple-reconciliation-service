package _helper

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"simple-reconciliation-service/cmd/root"
	"strconv"
	"strings"
)

func InitProgressBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSpinnerType(17),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func InitCommonArgs(extraArgs [][]string) [][]string {
	formatText := "-%s --%s"
	args := [][]string{
		{
			fmt.Sprintf(formatText, root.FlagFromDateShort, root.FlagFromDate),
			root.FlagFromDateValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagToDateShort, root.FlagToDate),
			root.FlagToDateValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagSystemTRXPathShort, root.FlagSystemTRXPath),
			root.FlagSystemTRXPathValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagBankTRXPathShort, root.FlagBankTRXPath),
			root.FlagBankTRXPathValue,
		},
		{
			fmt.Sprintf(formatText, root.FlagListBankShort, root.FlagListBank),
			strings.Join(root.FlagListBankValue, ","),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsVerboseShort, root.FlagIsVerbose),
			strconv.FormatBool(root.FlagIsVerboseValue),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsDebugShort, root.FlagIsDebug),
			strconv.FormatBool(root.FlagIsDebugValue),
		},
	}

	args = append(args, extraArgs...)
	return args
}
