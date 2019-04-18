package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start [name of the container]",
	Short: "Start a ipfs container",
	Args: cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		names := GetRunningNames()
		if contains(names, name) {
			fmt.Printf("'%s' is already running\n", name)
			return
		}

		execCmd("docker", false, "start", name)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
