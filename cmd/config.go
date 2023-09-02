package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xzzpig/simple-ip-kvm/internal/config"
	encoder "github.com/zwgblue/yaml-encoder"
	"go.uber.org/zap"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Config",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func printConfig(cfg *config.Config) {
	res, err := encoder.NewEncoder(cfg, encoder.WithComments(encoder.CommentsOnHead)).Encode()
	if err != nil {
		zap.L().Error("Fail to encode config to yaml", zap.Error(err))
		return
	}
	_, err = os.Stdout.Write(res)
	if err != nil {
		zap.L().Error("Fail to echo config", zap.Error(err))
		return
	}
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print current config",
	Run: func(cmd *cobra.Command, args []string) {
		printConfig(config.GetConfig())
	},
}

var whichCmd = &cobra.Command{
	Use:   "where",
	Short: "Which config file current using",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.ConfigFileUsed())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(printCmd)
	configCmd.AddCommand(whichCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
