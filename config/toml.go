package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var defaultConfigTmpl = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

ipfs_docker = "~/ipfs_docker"
`

// EnsureRoot ensure the root dir for ipd exists.
func EnsureRoot(rootDir string) {
	ensureDir(rootDir, 0700)

	cfgFilePath := path.Join(rootDir, "config.toml")
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		err := ioutil.WriteFile(cfgFilePath, []byte(defaultConfigTmpl), 0700)
		if err != nil {
			fmt.Printf("initialize config file error. %v\n", err)
			os.Exit(1)
		}
	}
}

func ensureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, mode); err != nil {
			return fmt.Errorf("could not create directory %v. %v", dir, err)
		}
	}

	return nil
}
