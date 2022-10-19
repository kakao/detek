package cmd

import (
	"fmt"
	"os"

	"github.com/kakao/detek/pkg/log"
	"github.com/kakao/detek/pkg/renderer"
	"github.com/spf13/cobra"
)

var (
	IsDebug bool
)

var rootCmd = &cobra.Command{
	Use:   "detek",
	Short: "Detecting Kubernetes known issues",
	Long:  `Detect is a cluster diagnostic tool, which aims to detect known issues automatically.`,
}

func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(func() {
		outputFormat = renderer.Format(outputFormstS)
		if err := outputFormat.IsValid(); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	})

	rootCmd.PersistentFlags().BoolVarP(&IsDebug, "debug", "d", false, "print logs for debugging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if IsDebug {
		log.EnableDebugMode()
	}
}
