package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var clearCmd = &cobra.Command{
	Use: "clear",
	Short: "Clear ipfs container",
	Run: func(cmd *cobra.Command, args []string) {
		running := GetRunningNames()
		for _, n := range running {
			stopContainer(n)
		}
		exited := getExitedNames()
		for _, n := range exited {
			removeContainer(n)
		}

		ipfsDir := getIpfsDir()

		curFile := path.Join(ipfsDir, currentName)
		os.Remove(curFile)

		cidFile := path.Join(ipfsDir, cidName)
		os.Remove(cidFile)

		tmpFolder := path.Join(ipfsDir, "tmp")
		os.RemoveAll(tmpFolder)

		fmt.Println("All ipfs containers have been cleared")
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
