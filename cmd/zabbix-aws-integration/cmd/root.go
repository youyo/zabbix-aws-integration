package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	ZabbixAwsIntegration struct {
		AccessKey       string
		SecretAccessKey string
		Arn             string
		Region          string
		Debug           bool
		Ec2             struct {
			Discovery struct {
				ZabbixHostGroup string
			}
		}
	}
)

var (
	Name                 string
	Version              string
	CommitHash           string
	BuildTime            string
	GoVersion            string
	cfgFile              string
	zabbixAwsIntegration ZabbixAwsIntegration
)

var RootCmd = &cobra.Command{
	Use: "zabbix-aws-integration",
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
	RootCmd.PersistentFlags().StringVarP(&zabbixAwsIntegration.AccessKey, "access-key", "", "", "AccessKey")
	RootCmd.PersistentFlags().StringVarP(&zabbixAwsIntegration.SecretAccessKey, "secret-access-key", "", "", "SecretAccessKey")
	RootCmd.PersistentFlags().StringVarP(&zabbixAwsIntegration.Arn, "arn", "", "", "Arn")
	RootCmd.PersistentFlags().StringVarP(&zabbixAwsIntegration.Region, "region", "", "ap-northeast-1", "Region")
	RootCmd.PersistentFlags().BoolVarP(&zabbixAwsIntegration.Debug, "debug", "", false, "Debug")
}

func initConfig() {
	viper.SetEnvPrefix("ai")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
