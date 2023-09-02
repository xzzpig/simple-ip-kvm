package config

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".simple-ip-kvm" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".simple-ip-kvm")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

	initConfigDefault()

	viper.Unmarshal(&cfg)

	InitZap()

	checkConfig()
}

func initConfigDefault() {
	viper.SetDefault("debug", cfg.Debug)
	viper.SetDefault("web.addr", "0.0.0.0:8000")
	viper.SetDefault("web.basepath", "/")
	viper.SetDefault("web.title", "Simple IP-KVM")
	viper.SetDefault("serialport", "/dev/ttyUSB0")
	viper.SetDefault("video.type", StreamerTypeExternal)
	viper.SetDefault("video.StreamUrl", "{URL_PROTOCOL}//{URL_HOST}:{STREAM_PORT}")
	viper.SetDefault("video.StreamPath", "/?action=stream")
	viper.SetDefault("video.SnapshotPath", "/?action=stream")
	viper.SetDefault("video.streamer.port", 8010)
	viper.SetDefault("proxy.enable", false)
	viper.SetDefault("proxy.path", "/stream")
	viper.SetDefault("proxy.rewrite", true)
}

func checkConfig() {
	cfg.Proxy.Path = strings.TrimRight(cfg.Proxy.Path, "/")
}
