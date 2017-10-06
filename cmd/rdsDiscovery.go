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
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"github.com/youyo/zabbix-userparameter-script-aws-integration/lib/aws-integration"
)

func integratedRdsDiscovery() (d aws_integration.DiscoveryData, err error) {
	if awsIntegrated.Arn != "" {
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, awsIntegrated.Arn)
		integratedRds.Service = rds.New(
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
	params := &rds.DescribeDBInstancesInput{}
	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, 3*time.Second)
	defer cancelFn()
	resp, err := integratedRds.Service.DescribeDBInstancesWithContext(ctx, params)
	if err != nil {
		return
	}
	for _, dbInstance := range resp.DBInstances {
		d = append(d, aws_integration.DiscoveryItem{
			"DB_INSTANCE_IDENTIFIER": *dbInstance.DBInstanceIdentifier,
			"DB_INSTANCE_ARN":        *dbInstance.DBInstanceArn,
			"DB_INSTANCE_CLASS":      *dbInstance.DBInstanceClass,
			"DB_INSTANCE_STATUS":     *dbInstance.DBInstanceStatus,
			"ENGINE":                 *dbInstance.Engine,
			"ENGINE_VERSION":         *dbInstance.EngineVersion,
		})
	}
	return
}

var rdsDiscoveryCmd = &cobra.Command{
	Use: "discovery",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := integratedRdsDiscovery()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		fmt.Println(d.Json())
	},
}

func init() {
	rdsCmd.AddCommand(rdsDiscoveryCmd)
}
