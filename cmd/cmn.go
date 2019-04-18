package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	currentName = "current"
	cidName     = "cid"
)

func getIpfsDir() string {
	ipfsDir := viper.GetString("ipfs_docker")
	ipfsDir, err := homedir.Expand(ipfsDir)
	if err != nil {
		fmt.Printf("ipfs_docker setting is invalid. %v\n", err)
		os.Exit(1)
	}

	return ipfsDir
}

func execCmd(name string, retOut bool, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if retOut {
		return out.String()
	} else {
		fmt.Println(out.String())
		return ""
	}
}

func contains(names []string, name string) bool {
	if len(names) == 0 {
		return false
	}

	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func rmEmpty(elems []string) []string {
	for i := 0; i < len(elems); i++ {
		if elems[i] == "" {
			elems = append(elems[:i], elems[i+1:]...)
		}
	}
	return elems
}

func getStaging(name string) string {
	ind := strings.LastIndex(name, "_")
	num := name[ind+1:]

	ipfsDir := getIpfsDir()
	ipfsStaging := path.Join(ipfsDir, "tmp", fmt.Sprintf("ipfs_staging_%s", num))
	return ipfsStaging
}

func writeCID(cid string) {
	fPath := path.Join(getIpfsDir(), cidName)
	file, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Write cid failed")
		os.Exit(1)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	w.WriteString(fmt.Sprintf("%s\n", cid))
	w.Flush()
}

func getCIDs() []string {
	fPath := path.Join(getIpfsDir(), cidName)
	bz, err := ioutil.ReadFile(fPath)
	if os.IsNotExist(err) {
		return []string{}
	}

	if err != nil {
		fmt.Println("Read cid failed")
		os.Exit(1)
	}

	lines := strings.Split(string(bz), "\n")
	cids := rmEmpty(lines)
	return cids
}

func pinAdd(name, cid string) {
	execCmd("docker", false, "exec", name, "ipfs", "pin", "add", cid)
}
