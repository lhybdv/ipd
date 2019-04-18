package cmd

import (
	"fmt"
	"github.com/lhybdv/ipd/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var rootCmd = &cobra.Command{
	Use:   "ipd",
	Short: "Ipd is a assistant for using ipfs with docker",
}

var cfgFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfgPath := path.Join(home, ".ipd")
	viper.SetConfigType("toml")
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		config.EnsureRoot(cfgPath)
		viper.ReadInConfig()
	}
}
