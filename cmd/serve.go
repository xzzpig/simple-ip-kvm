package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xzzpig/simple-ip-kvm/internal/hid"
	"github.com/xzzpig/simple-ip-kvm/internal/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch ip kvm server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		hid.SetupVideoSerial()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return hid.CLose()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
