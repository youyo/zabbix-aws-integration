package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

func integratedCloudWatchGetMetrics() (value float64, unit string, err error) {
	if awsIntegrated.Arn != "" {
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, awsIntegrated.Arn)
		integratedCloudWatch.Service = cloudwatch.New(
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
	layout := "2006-01-02T15:04:05Z"
	endtime, _ := time.Parse(layout, time.Now().UTC().Format(layout))
	starttime := endtime.Add(-600 * time.Second)
	params := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(integratedCloudWatch.DimensionName),
				Value: aws.String(integratedCloudWatch.DimensionValue),
			},
		},
		Namespace:  aws.String(integratedCloudWatch.Namespace),
		MetricName: aws.String(integratedCloudWatch.MetricName),
		Period:     aws.Int64(60),
		EndTime:    &endtime,
		StartTime:  &starttime,
		Statistics: []*string{
			aws.String("Minimum"),
			aws.String("Maximum"),
			aws.String("Average"),
			aws.String("SampleCount"),
			aws.String("Sum"),
		},
	}
	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, 3*time.Second)
	defer cancelFn()
	resp, err := integratedCloudWatch.Service.GetMetricStatisticsWithContext(ctx, params)
	if err != nil {
		return
	}
	if len(resp.Datapoints) > 0 {
		sort.Slice(resp.Datapoints, func(i, j int) bool {
			return resp.Datapoints[i].Timestamp.Unix() > resp.Datapoints[j].Timestamp.Unix()
		})
		if awsIntegrated.Debug {
			pp.Print(resp)
		}
		switch integratedCloudWatch.Statistics {
		case "Minimum":
			value = *resp.Datapoints[0].Minimum
		case "Maximum":
			value = *resp.Datapoints[0].Maximum
		case "Average":
			value = *resp.Datapoints[0].Average
		case "SampleCount":
			value = *resp.Datapoints[0].SampleCount
		case "Sum":
			value = *resp.Datapoints[0].Sum
		default:
			err = errors.New("Statistics is not match")
		}
		unit = *resp.Datapoints[0].Unit
	}
	return
}

var cloudWatchGetMetricsCmd = &cobra.Command{
	Use:   "get-metrics",
	Short: "A brief description of your command",
	Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		value, unit, err := integratedCloudWatchGetMetrics()
		if err != nil {
			pp.Print(err)
			os.Exit(1)
		}
		if integratedCloudWatch.Unit {
			fmt.Println(unit)
		} else {
			fmt.Println(value)
		}
	},
}

func init() {
	cloudwatchCmd.AddCommand(cloudWatchGetMetricsCmd)
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&integratedCloudWatch.DimensionName, "dimention-name", "", "", "DimensionName")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&integratedCloudWatch.DimensionValue, "dimention-value", "", "", "DimensionValue")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&integratedCloudWatch.Namespace, "namespace", "", "", "Namespace")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&integratedCloudWatch.MetricName, "metric-name", "", "CPUUtilization", "MetricName")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&integratedCloudWatch.Statistics, "statistics", "", "Average", "Statistics")
	cloudWatchGetMetricsCmd.PersistentFlags().BoolVarP(&integratedCloudWatch.Unit, "unit", "", false, "Unit")
}
