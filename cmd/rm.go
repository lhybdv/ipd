package cmd

import "github.com/spf13/cobra"

var (
	name  string
	rmCmd = &cobra.Command{
		Use:   "rm [name of container]",
		Short: "Remove container by name",
		Run: func(cmd *cobra.Command, args []string) {
			if contains(GetRunningNames(), name) {
				stopContainer(name)
			} else if !contains(getExitedNames(), name) {
				return
			}

			removeContainer(name)
		},
	}
)

func init() {
	addRmArgs(rmCmd)
	rootCmd.AddCommand(rmCmd)
}

func addRmArgs(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&name, "name", "n", "", "name of container")
}

func removeContainer(name string) {
	execCmd("docker", false, "rm", name)
}
