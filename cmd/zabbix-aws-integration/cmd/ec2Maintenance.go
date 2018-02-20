package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var (
	ec2InstanceID        string
	noMaintenanceMessage string
)

func fetchInstanceStatus(ec2Service *ec2.EC2, ec2InstanceID string) (resp *ec2.DescribeInstanceStatusOutput, err error) {
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&ec2InstanceID},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = ec2Service.DescribeInstanceStatusWithContext(ctx, params)
	return
}

func buildMaintenanceMessage(resp *ec2.DescribeInstanceStatusOutput, noMaintenanceMessage string) (message string) {
	message = noMaintenanceMessage
	if len(resp.InstanceStatuses[0].Events) > 0 {
		message = fmt.Sprintf("Code: %s, Description: %s, NotAfter: %s, NotBefore: %s",
			*resp.InstanceStatuses[0].Events[0].Code,
			*resp.InstanceStatuses[0].Events[0].Description,
			*resp.InstanceStatuses[0].Events[0].NotAfter,
			*resp.InstanceStatuses[0].Events[0].NotBefore,
		)
	}
	return
}

func ec2Maintenance(arn, accessKey, secretAccessKey, region, ec2InstanceID, noMaintenanceMessage string) (message string, err error) {
	sess, creds, err := NewCredentials(
		arn,
		accessKey,
		secretAccessKey,
	)
	if err != nil {
		return
	}
	config := NewConfig(creds, region)
	ec2Service := NewEc2Service(sess, config)
	resp, err := fetchInstanceStatus(ec2Service, ec2InstanceID)
	if err != nil {
		return
	}
	message = buildMaintenanceMessage(resp, noMaintenanceMessage)
	return
}

var ec2MaintenanceCmd = &cobra.Command{
	Use: "maintenance",
	//Short: "",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		message, err := ec2Maintenance(
			zabbixAwsIntegration.Arn,
			zabbixAwsIntegration.AccessKey,
			zabbixAwsIntegration.SecretAccessKey,
			zabbixAwsIntegration.Region,
			ec2InstanceID,
			noMaintenanceMessage,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2MaintenanceCmd)
	ec2MaintenanceCmd.PersistentFlags().StringVarP(&ec2InstanceID, "instance-id", "i", "", "InstanceID")
	ec2MaintenanceCmd.PersistentFlags().StringVarP(&noMaintenanceMessage, "no-maintenance-message", "m", "There is no maintenance", "No Maintenance Message")
}
