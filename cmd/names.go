package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var (
	all      bool
	exited   bool
	namesCmd = &cobra.Command{
		Use:   "names",
		Short: "Get ipfs container names from docker",
		Run: func(cmd *cobra.Command, args []string) {
			var status string
			var names []string
			if all {
				names = getAllNames()
			} else if exited {
				status = "exited"
				names = getExitedNames()
			} else {
				status = "running"
				names = GetRunningNames()
			}

			if len(names) == 0 {
				if status == "" {
					fmt.Println("No container")
				} else {
					fmt.Printf("No %s container\n", status)
				}
				return
			}

			for _, name := range names {
				fmt.Println(name)
			}
		},
	}
)

func init() {
	addNamesFlags(namesCmd)
	rootCmd.AddCommand(namesCmd)
}

func addNamesFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers")
	cmd.Flags().BoolVarP(&exited, "exited", "e", false, "Show exited containers")
}

func getAllNames() []string {
	running := GetRunningNames()
	exited := getExitedNames()
	all := append(running, exited...)
	return all
}

func GetRunningNames() []string {
	return getNames("running")
}

func getExitedNames() []string {
	return getNames("exited")
}

func getNames(status string) []string {
	args := []string{
		"ps",
		"--format", "{{.Names}}",
		"--filter", "name=ipfs",
		"--filter", fmt.Sprintf("status=%s", status),
	}
	out := execCmd("docker", true, args...)
	if len(out) == 0 {
		return []string{}
	}

	names := strings.Split(out, "\n")
	names = rmEmpty(names)
	return names
}
