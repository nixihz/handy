package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"handy/utils"
	"os"
)

var (
	// Used for configs file.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "handy",
		Short: "a handy tool box for you, lao tie!",
		Long:  `a handy tool box for you, lao tie!`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "configs", "c", "prod", "configs file (default is $HOME/configs/prod.yaml)")
}

func initConfig() {
	filePath := utils.GetRunPath() + "/configs"
	viper.SetConfigName(cfgFile)
	viper.AutomaticEnv()
	viper.AddConfigPath(filePath)
	workDir := os.Getenv("HANDY_WORK_DIR")
	if workDir == "" {
		fmt.Print("error HANDY_WORK_DIR env is not set.")
		os.Exit(2)
	}

	viper.AddConfigPath(workDir + "/configs")
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("build/configs")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
	}
}
