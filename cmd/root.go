package cmd

import (
	"encoding/json"
	"examresult/config"
	"examresult/model"
	"examresult/server"
	"examresult/service"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "examresult",
	Short: "exam result notice server",
	Long:  `examresult exam result notice server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init
		config.Init(cfgFile)
		model.Init()
		server.Init()
		service.Init()

		// Run
		service.Run()
		server.Run() // gin run 在这个里面，必须最后启
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	configExample, _ := json.Marshal(config.Conf{})

	rootCmd.Flags().StringVarP(&cfgFile,
		"config", "c", "",
		fmt.Sprintf("config file, e.g. %s", configExample))
	_ = rootCmd.MarkFlagRequired("config")
}
