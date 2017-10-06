package cmd

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

type (
	IntegratedRds struct {
		Action            string
		Service           *rds.RDS
		DbInstanceID      string
		CloudWatchService *cloudwatch.CloudWatch
	}
)

var integratedRds IntegratedRds

var rdsCmd = &cobra.Command{
	Use: "rds",
	//Short: "A brief description of your command",
	//Long: `A brief description of your command`,
}

func init() {
	RootCmd.AddCommand(rdsCmd)
	//rdsCmd.Flags().StringVarP(&integratedRds., "", "", "", "")
}
