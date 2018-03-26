package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type (
	Ec2DiscoveryItem struct {
		InstanceID        string `json:"{#INSTANCE_ID}"`
		InstanceName      string `json:"{#INSTANCE_NAME}"`
		InstanceRole      string `json:"{#INSTANCE_ROLE}"`
		InstancePublicIp  string `json:"{#INSTANCE_PUBLIC_IP}"`
		InstancePrivateIp string `json:"{#INSTANCE_PRIVATE_IP}"`
		IfConn            string `json:"{#IF.CONN}"`
		IfIp              string `json:"{#IF.IP}"`
		IfDns             string `json:"{#IF.DNS}"`
		IfPort            string `json:"{#IF.PORT}"`
		IfType            string `json:"{#IF.TYPE}"`
		IfDefault         int    `json:"{#IF.DEFAULT}"`
	}
	Ec2DiscoveryItems []Ec2DiscoveryItem
	Ec2DiscoveryData  struct {
		Data Ec2DiscoveryItems `json:"data"`
	}
)

func fetchRunningInstances(ec2Service *ec2.EC2) (resp *ec2.DescribeInstancesOutput, err error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = ec2Service.DescribeInstancesWithContext(ctx, params)
	return
}

func fetchInstanceName(ec2Instance *ec2.Instance) (instnaceName string) {
	for _, tag := range ec2Instance.Tags {
		if *tag.Key == "Name" {
			instnaceName = *tag.Value
		}
	}
	return
}

func fetchInstanceRole(ec2Instance *ec2.Instance) (instnaceRole string) {
	for _, tag := range ec2Instance.Tags {
		if *tag.Key == "Role" {
			instnaceRole = *tag.Value
		}
	}
	return
}

func buildEc2DiscoveryData(resp *ec2.DescribeInstancesOutput) (ec2DiscoveryData Ec2DiscoveryData, err error) {
	var ec2DiscoveryItems Ec2DiscoveryItems
	for _, v := range resp.Reservations {
		for _, i := range v.Instances {
			instanceName := fetchInstanceName(i)
			instanceRole := fetchInstanceRole(i)
			ec2DiscoveryItems = append(ec2DiscoveryItems, Ec2DiscoveryItem{
				InstanceID:        *i.InstanceId,
				InstanceName:      instanceName,
				InstanceRole:      instanceRole,
				InstancePublicIp:  *i.PublicIpAddress,
				InstancePrivateIp: *i.PrivateIpAddress,
				IfConn:            *i.PublicIpAddress,
				IfIp:              *i.PublicIpAddress,
				IfDns:             instanceName,
				IfPort:            "10050",
				IfType:            "AGENT",
				IfDefault:         1,
			})
		}
	}
	ec2DiscoveryData = Ec2DiscoveryData{ec2DiscoveryItems}
	return
}

func ec2Discovery(arn, accessKey, secretAccessKey, region string) (ec2DiscoveryData Ec2DiscoveryData, err error) {
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
	resp, err := fetchRunningInstances(ec2Service)
	if err != nil {
		return
	}
	ec2DiscoveryData, err = buildEc2DiscoveryData(resp)
	return
}

var ec2DiscoveryCmd = &cobra.Command{
	Use: "discovery",
	//Short: "",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := ec2Discovery(
			zabbixAwsIntegration.Arn,
			zabbixAwsIntegration.AccessKey,
			zabbixAwsIntegration.SecretAccessKey,
			zabbixAwsIntegration.Region,
		)
		if err != nil {
			log.Fatal(err)
		}
		s, err := Jsonize(d)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2DiscoveryCmd)
}
