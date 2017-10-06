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

func integratedRdsClusterDiscovery() (d aws_integration.DiscoveryData, err error) {
	if awsIntegrated.Arn != "" {
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, awsIntegrated.Arn)
		integratedRdsCluster.Service = rds.New(
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
	params := &rds.DescribeDBClustersInput{}
	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, 3*time.Second)
	defer cancelFn()
	resp, err := integratedRdsCluster.Service.DescribeDBClustersWithContext(ctx, params)
	if err != nil {
		return
	}
	for _, dbCluster := range resp.DBClusters {
		d = append(d, aws_integration.DiscoveryItem{
			"DB_CLUSTER_IDENTIFIER": *dbCluster.DBClusterIdentifier,
			"DB_CLUSTER_ARN":        *dbCluster.DBClusterArn,
			"STATUS":                *dbCluster.Status,
			"ENGINE":                *dbCluster.Engine,
			"ENGINE_VERSION":        *dbCluster.EngineVersion,
		})
	}
	return
}

var rdsClusterDiscoveryCmd = &cobra.Command{
	Use: "discovery",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := integratedRdsClusterDiscovery()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		fmt.Println(d.Json())
	},
}

func init() {
	rdsClusterCmd.AddCommand(rdsClusterDiscoveryCmd)
}
