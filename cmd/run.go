package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	maxCount    = 5
	maxTotal    = 10
	keyFileName = "swarm.key"
)

var runCmd = &cobra.Command{
	Use:   "run [count of ipfs instances]",
	Short: "Run one or more ipfs instances",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cntStr := args[0]
		cnt64, err := strconv.ParseInt(cntStr, 10, 32)
		if err != nil {
			fmt.Println("[count] should be a integer")
			return
		}

		cnt := int(cnt64)
		if cnt > maxCount {
			fmt.Printf("The maxmium count is %d\n", maxCount)
			return
		}

		index := getCurent()
		if index+cnt > maxTotal {
			fmt.Printf("The maxmium total is %d\n", maxTotal)
			return
		}

		var names []string
		go func() {
			for i := index + 1; i <= index+cnt; i++ {
				name := run(i)
				names = append(names, name)
			}
		}()

		// Wait for ipfs container initializing
		time.Sleep(time.Second * 15)

		cids := getCIDs()
		for _, n := range names {
			for _, cid := range cids {
				pinAdd(n, cid)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(index int) string {
	ipfsDir := getIpfsDir()

	ipfsPaging := path.Join(ipfsDir, "tmp", fmt.Sprintf("ipfs_staging_%d", index))
	renewDir(ipfsPaging)
	ipfsData := path.Join(ipfsDir, "tmp", fmt.Sprintf("ipfs_data_%d", index))
	renewDir(ipfsData)

	keyDir := path.Join(path.Join(ipfsDir, "tmp", fmt.Sprintf("ipfs_data_%d", index)), keyFileName)
	err := copyKey(ipfsDir, keyDir)
	if err != nil {
		fmt.Printf("copy key failed. %v\n", err)
	}

	hostName := fmt.Sprintf("ipfs_host_%d", index)
	fmt.Printf("Creating %s\n", hostName)

	argsRun := buildArgsRun(hostName, ipfsPaging, ipfsData, index)
	execCmd("docker", false, argsRun...)

	//fmt.Printf("Remove bootstrap for %s", hostName)

	//argsRemove := buildArgsRemove(hostName)
	//execCmd("docker", false, argsRemove...)

	setCurrent(index)
	return hostName
}

func renewDir(dir string) {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0700)
}

func copyKey(ipfsDir, destDir string) error {
	srcDir := path.Join(ipfsDir, "data", keyFileName)
	src, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destDir)
	if err != nil {
		return err
	}
	defer dest.Close()

	io.Copy(dest, src)

	return nil
}

func buildArgsRun(hostName, ipfsStaging, ipfsData string, index int) []string {
	return []string{
		"run", "-d",
		"--name", hostName,
		"-v", fmt.Sprintf("%s:/export", ipfsStaging),
		"-v", fmt.Sprintf("%s:/data/ipfs", ipfsData),
		"-p", fmt.Sprintf("%d:4001", 4000+index),
		"-p", fmt.Sprintf("%d:5001", 5000+index),
		"-p", fmt.Sprintf("127.0.0.1:%d:8080", 8080+index),
		"ipfs/go-ipfs:latest",
	}
}

func buildArgsRemove(hostName string) []string {
	return []string{"exec", hostName, "ipfs", "bootstrap", "rm", "--all"}
}

func getCurent() int {
	ipfsDir := getIpfsDir()

	curPath := path.Join(ipfsDir, currentName)

	if _, err := os.Stat(curPath); os.IsNotExist(err) {
		writeCurrent(curPath, 0)
		return 0
	}

	bz, err := ioutil.ReadFile(curPath)
	if err != nil {
		writeCurrent(curPath, 0)
		return 0
	}

	index, err := strconv.ParseInt(string(bz), 10, 32)
	if err != nil {
		writeCurrent(curPath, 0)
		return 0
	}

	return int(index)
}

func setCurrent(index int) {
	curPath := path.Join(getIpfsDir(), currentName)
	writeCurrent(curPath, index)
}

func writeCurrent(curPath string, index int) {
	file, err := os.OpenFile(curPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("write current file failed")
		os.Exit(1)
	}
	defer file.Close()

	indexStr := strconv.Itoa(index)
	w := bufio.NewWriter(file)
	w.WriteString(indexStr)
	w.Flush()
}
