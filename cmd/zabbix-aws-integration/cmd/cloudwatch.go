package cmd

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type (
	IntegratedCloudWatch struct {
		Service        *cloudwatch.CloudWatch
		DimensionName  string
		DimensionValue string
		Namespace      string
		MetricName     string
		Statistics     string
		Unit           bool
	}
)

var integratedCloudWatch IntegratedCloudWatch

var cloudwatchCmd = &cobra.Command{
	Use: "cloudwatch",
	//Short: "A brief description of your command",
	//Long: `A brief description of your command`,
}

func init() {
	RootCmd.AddCommand(cloudwatchCmd)
}
