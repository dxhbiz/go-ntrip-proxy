package server

import (
	"fmt"
	"os"
	"path"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/config"
	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/exe"
	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/log"
	"github.com/dxhbiz/go-ntrip-proxy/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	APP_NAME = "ntrip-proxy"
)

var (
	// the path of the configuration file
	cfgFile = ""
	// the path of the application
	exePath = exe.Path()
)

var rootCmd = &cobra.Command{
	Use:   APP_NAME,
	Short: "An ntrip proxy forwarding tool",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $EXE/config/config.json)")
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	var configFile string
	if cfgFile != "" {
		if path.IsAbs(cfgFile) {
			configFile = cfgFile
		} else {
			configFile = path.Join(exePath, cfgFile)
		}
	} else {
		configFile = path.Join(exePath, "config", "config.json")
	}
	viper.SetConfigFile(configFile)

	err := config.InitConfig(configFile)
	if err != nil {
		fmt.Printf("Init config file %s error: %s", configFile, err.Error())
		os.Exit(1)
	}
}

// Execute start the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initConfig()

	cfg := config.GetConfig()

	log.Init(cfg.Log)
	defer log.Sync()

	log.Infof("%s version: %s", APP_NAME, version.RELEASE)
}
