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
	"github.com/youyo/zabbix-userparameter-script-aws-integration/lib/aws-integration"
)

func integratedEc2Discovery() (d aws_integration.DiscoveryData, err error) {
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
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, 3*time.Second)
	defer cancelFn()
	resp, err := integratedEc2.Service.DescribeInstancesWithContext(ctx, params)
	if err != nil {
		return
	}
	for _, v := range resp.Reservations {
		for _, i := range v.Instances {
			integratedEc2.InstanceID = *i.InstanceId
			instanceName := func() string {
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						return *t.Value
					}
				}
				return ""
			}()
			d = append(d, aws_integration.DiscoveryItem{
				"INSTANCE_ID":   integratedEc2.InstanceID,
				"INSTANCE_NAME": instanceName,
			})
		}
	}
	return
}

var ec2DiscoveryCmd = &cobra.Command{
	Use: "discovery",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := integratedEc2Discovery()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		fmt.Println(d.Json())
	},
}

func init() {
	ec2Cmd.AddCommand(ec2DiscoveryCmd)
}
