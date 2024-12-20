package process

import (
	"simple-reconciliation-service/cmd/_helper"
	"simple-reconciliation-service/cmd/root"
	"simple-reconciliation-service/internal/app/component/csqlite"

	"github.com/spf13/cobra"
)

func Runner(cmd *cobra.Command, args []string) (er error) {
	dBPath := csqlite.DBPath{}

	if root.FlagIsDebugValue {
		dBPath.WriteDBPath = "./reconciliation.db"
	}

	return _helper.RunnerSubCommand(cmd, args, dBPath)
}
