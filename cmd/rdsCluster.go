package cmd

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

type (
	IntegratedRdsCluster struct {
		Action            string
		Service           *rds.RDS
		DbClusterID       string
		CloudWatchService *cloudwatch.CloudWatch
	}
)

var integratedRdsCluster IntegratedRdsCluster

var rdsClusterCmd = &cobra.Command{
	Use: "rds-cluster",
	//Short: "A brief description of your command",
	//Long: `A brief description of your command`,
}

func init() {
	RootCmd.AddCommand(rdsClusterCmd)
	//RootCmd.Flags().StringVarP(&integratedRds., "", "", "", "")
}
