package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"path"
)

var getCmd = &cobra.Command{
	Use: "get [name of ipfs container] [cid] [directory]",
	Short: "Download ipfs objects by [cid] to [direcotry]",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cid := args[1]
		dir := args[2]
		err := Get(name,cid ,dir)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func Get(name, cid, dir string) error {
	execCmd("docker", false, "exec", name, "ipfs", "get", cid, "-o", "/export")
	dir, err := homedir.Expand(dir)
	if err != nil {
		return err
	}
	ipfsStaging := getStaging(name)

	fPath := path.Join(ipfsStaging, cid)

	execCmd("cp", false,  "-r", fPath, dir)

	execCmd("rm", false, "-r", fPath)
	return nil
}
