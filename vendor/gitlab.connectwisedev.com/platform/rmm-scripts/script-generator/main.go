package main

import (
	"gitlab.connectwisedev.com/RMM/rmm-scripts/script-generator/cmd"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "script-generator"}

	rootCmd.AddCommand(cmd.Create)
	rootCmd.AddCommand(cmd.Update)
	rootCmd.Execute()
}
