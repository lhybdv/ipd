package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"strings"
)

var addCmd = &cobra.Command{
	Use: "add [name of ipfs container] [file or folder]",
	Short: "Add file or folder to ipfs node",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fPath := args[1]
		cid := Add(name, fPath)
		fmt.Println(cid)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func Add(name, fPath string) string {
	names := GetRunningNames()
	if !contains(names, name) {
		fmt.Printf("'%s' is not running", name)
		return ""
	}

	ipfsStaging := getStaging(name)

	fPath, err := homedir.Expand(fPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	execCmd("cp", false, "-r", fPath, ipfsStaging)
	fileName := filepath.Base(fPath)
	out := execCmd("docker", true, "exec", name, "ipfs", "add", "-r", fmt.Sprintf("/export/%s", fileName))
	execCmd("rm", false, "-r", path.Join(ipfsStaging, fileName))
	arr := strings.Split(out, " ")
	if len(arr) == 3 {
		cid := arr[1]
		writeCID(cid)

		names := GetRunningNames()
		for _, n := range names {
			if n != name {
				pinAdd(n, cid)
			}
		}
		return cid
	}
	fmt.Println("Some error for adding")
	return ""
}

