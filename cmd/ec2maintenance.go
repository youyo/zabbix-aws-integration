package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

func integratedEc2Maintenance() (value string, err error) {
	value = "There is no maintenance"
	if awsIntegrated.Arn != "" {
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, awsIntegrated.Arn)
		integratedEc2.Service = ec2.New(
			sess,
			aws.NewConfig().WithRegion(awsIntegrated.Region).WithCredentials(creds),
		)
	} else if awsIntegrated.AccessKey != "" && awsIntegrated.SecretAccessKey != "" {
		err = errors.New("not yet implemented")
		return
	} else {
		err = errors.New("credentials are not enough")
		return
	}
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&integratedEc2.InstanceID},
	}
	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, 3*time.Second)
	defer cancelFn()
	resp, err := integratedEc2.Service.DescribeInstanceStatusWithContext(ctx, params)
	if err != nil {
		return
	}
	if len(resp.InstanceStatuses[0].Events) > 0 {
		value = fmt.Sprintf("Code: %s, Description: %s, NotAfter: %s, NotBefore: %s",
			*resp.InstanceStatuses[0].Events[0].Code,
			*resp.InstanceStatuses[0].Events[0].Description,
			*resp.InstanceStatuses[0].Events[0].NotAfter,
			*resp.InstanceStatuses[0].Events[0].NotBefore,
		)
	}
	return
}

var ec2MaintenanceCmd = &cobra.Command{
	Use: "maintenance",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := integratedEc2Maintenance()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		fmt.Println(v)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2MaintenanceCmd)
	ec2MaintenanceCmd.PersistentFlags().StringVarP(&integratedEc2.InstanceID, "instance-id", "i", "", "InstanceID")
}
