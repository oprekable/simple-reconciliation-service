package sample

import (
	"github.com/spf13/cobra"
	"simple-reconciliation-service/cmd/_helper"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component/csqlite"
)

func Runner(cmd *cobra.Command, args []string) (er error) {
	dBPath := csqlite.DBPath{}

	if root.FlagIsDebugValue {
		dBPath.ReadDBPath = "./sample.db"
	}

	return _helper.RunnerSubCommand(cmd, args, dBPath)
}
