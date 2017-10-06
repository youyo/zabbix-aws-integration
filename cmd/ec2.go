package cmd

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type (
	IntegratedEc2 struct {
		Action            string
		Service           *ec2.EC2
		InstanceID        string
		CloudWatchService *cloudwatch.CloudWatch
	}
)

var integratedEc2 IntegratedEc2

var ec2Cmd = &cobra.Command{
	Use: "ec2",
	//Short: "A brief description of your command",
	//Long: `A brief description of your command`,
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
	//ec2Cmd.Flags().StringVarP(&integratedEc2., "", "", "", "")
}
