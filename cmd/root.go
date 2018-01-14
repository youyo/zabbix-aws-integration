package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const noMaintenanceMessage string = "There is no maintenance"

type (
	AwsIntegrated struct {
		AccessKey       string
		SecretAccessKey string
		Arn             string
		Region          string
		Debug           bool
	}
)

var (
	Name          string
	Version       string
	CommitHash    string
	BuildTime     string
	GoVersion     string
	cfgFile       string
	awsIntegrated AwsIntegrated
)

var RootCmd = &cobra.Command{
	Use: "aws-integration",
	//Short: "A brief description of your application",
	//Long: `A brief description of your application`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&awsIntegrated.AccessKey, "access-key", "", "", "AccessKey")
	RootCmd.PersistentFlags().StringVarP(&awsIntegrated.SecretAccessKey, "secret-access-key", "", "", "SecretAccessKey")
	RootCmd.PersistentFlags().StringVarP(&awsIntegrated.Arn, "arn", "", "", "Arn")
	RootCmd.PersistentFlags().StringVarP(&awsIntegrated.Region, "region", "", "ap-northeast-1", "Region")
	RootCmd.PersistentFlags().BoolVarP(&awsIntegrated.Debug, "debug", "", false, "Debug")
}

func initConfig() {
	viper.SetEnvPrefix("ai")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
