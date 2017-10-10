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
)

func integratedRdsClusterMaintenance() (value string, err error) {
	value = "There is no maintenance"
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
	params := &rds.DescribePendingMaintenanceActionsInput{
		Filters: []*rds.Filter{
			{
				Name:   aws.String("db-cluster-id"),
				Values: []*string{aws.String(integratedRdsCluster.DbClusterID)},
			},
		},
	}
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFn()
	resp, err := integratedRdsCluster.Service.DescribePendingMaintenanceActionsWithContext(ctx, params)
	if err != nil {
		return
	}
	if len(resp.PendingMaintenanceActions) > 0 {
		value = fmt.Sprintf(
			"Action: %s, Description: %s",
			*resp.PendingMaintenanceActions[0].PendingMaintenanceActionDetails[0].Action,
			*resp.PendingMaintenanceActions[0].PendingMaintenanceActionDetails[0].Description,
		)
	}
	return
}

var rdsClusterMaintenanceCmd = &cobra.Command{
	Use: "maintenance",
	//Short: "A brief description of your command",
	//Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := integratedRdsClusterMaintenance()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		fmt.Println(v)
	},
}

func init() {
	rdsClusterCmd.AddCommand(rdsClusterMaintenanceCmd)
	rdsClusterMaintenanceCmd.PersistentFlags().StringVarP(&integratedRdsCluster.DbClusterID, "db-cluster-id", "", "", "DbClusterID")
}
