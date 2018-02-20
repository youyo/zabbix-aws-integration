package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type (
	CloudWatchCommandOptions struct {
		DimensionName  string
		DimensionValue string
		Namespace      string
		MetricName     string
		Statistics     string
		WithUnit       bool
	}
)

var cloudwatchCommandOptions CloudWatchCommandOptions

func buildRequestParams(dimensionName, dimensionValue, namespace, metricName string) (params *cloudwatch.GetMetricStatisticsInput, err error) {
	endTime, err := time.Parse(TimeLayout, time.Now().UTC().Format(TimeLayout))
	if err != nil {
		return
	}
	startTime := endTime.Add(-600 * time.Second)
	params = &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(60),
		EndTime:    &endTime,
		StartTime:  &startTime,
		Statistics: []*string{
			aws.String("Minimum"),
			aws.String("Maximum"),
			aws.String("Average"),
			aws.String("SampleCount"),
			aws.String("Sum"),
		},
	}
	return
}

func fetchCloudWatchMetrics(cloudWatchService *cloudwatch.CloudWatch, params *cloudwatch.GetMetricStatisticsInput) (resp *cloudwatch.GetMetricStatisticsOutput, err error) {
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = cloudWatchService.GetMetricStatisticsWithContext(ctx, params)
	return
}

func extractValues(resp *cloudwatch.GetMetricStatisticsOutput, statistics string) (value float64, unit string, err error) {
	if len(resp.Datapoints) > 0 {
		sort.Slice(resp.Datapoints, func(i, j int) bool {
			return resp.Datapoints[i].Timestamp.Unix() > resp.Datapoints[j].Timestamp.Unix()
		})
		switch statistics {
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
	} else {
		err = errors.New("Datapoint has not values")
	}
	return
}

func cloudWatchGetMetrics(arn, accessKey, secretAccessKey, region, dimensionName, dimensionValue, namespace, metricName, statistics string) (value float64, unit string, err error) {
	sess, creds, err := NewCredentials(
		arn,
		accessKey,
		secretAccessKey,
	)
	if err != nil {
		return
	}
	config := NewConfig(creds, region)
	cloudWatchService := NewCloudWatchService(sess, config)
	params, err := buildRequestParams(
		dimensionName,
		dimensionValue,
		namespace,
		metricName,
	)
	resp, err := fetchCloudWatchMetrics(cloudWatchService, params)
	if err != nil {
		return
	}
	value, unit, err = extractValues(resp, statistics)
	return
}

func outputValues(value float64, unit string, withUnit bool) {
	if withUnit {
		fmt.Fprintf(os.Stdout, "%f %s\n", value, unit)
	} else {
		fmt.Fprintf(os.Stdout, "%f\n", value)
	}
}

var cloudWatchGetMetricsCmd = &cobra.Command{
	Use:   "get-metrics",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		value, unit, err := cloudWatchGetMetrics(
			zabbixAwsIntegration.Arn,
			zabbixAwsIntegration.AccessKey,
			zabbixAwsIntegration.SecretAccessKey,
			zabbixAwsIntegration.Region,
			cloudwatchCommandOptions.DimensionName,
			cloudwatchCommandOptions.DimensionValue,
			cloudwatchCommandOptions.Namespace,
			cloudwatchCommandOptions.MetricName,
			cloudwatchCommandOptions.Statistics,
		)
		if err != nil {
			log.Fatal(err)
		}
		outputValues(value, unit, cloudwatchCommandOptions.WithUnit)
	},
}

func init() {
	cloudwatchCmd.AddCommand(cloudWatchGetMetricsCmd)
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&cloudwatchCommandOptions.DimensionName, "dimention-name", "", "InstanceId", "DimensionName")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&cloudwatchCommandOptions.DimensionValue, "dimention-value", "", "", "DimensionValue")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&cloudwatchCommandOptions.Namespace, "namespace", "", "AWS/EC2", "Namespace")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&cloudwatchCommandOptions.MetricName, "metric-name", "", "CPUUtilization", "MetricName")
	cloudWatchGetMetricsCmd.PersistentFlags().StringVarP(&cloudwatchCommandOptions.Statistics, "statistics", "", "Average", "Statistics")
	cloudWatchGetMetricsCmd.PersistentFlags().BoolVarP(&cloudwatchCommandOptions.WithUnit, "with-unit", "", false, "WithUnit")
}
