package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use: "stop [name of the container]",
	Short: "Stop the ipfs container",
	Args: cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		names := GetRunningNames()
		if !contains(names, name) {
			fmt.Printf("'%s' is not running\n", name)
			return
		}

		stopContainer(name)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopContainer(name string) {
	execCmd("docker", false, "stop", name)
}
