package cmd

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewCredentials(arn, accessKey, secretAccessKey string) (sess *session.Session, creds *credentials.Credentials, err error) {
	if arn != "" {
		sess = session.Must(session.NewSession())
		creds = stscreds.NewCredentials(sess, arn)
	} else if accessKey != "" && secretAccessKey != "" {
		err = errors.New("not yet implemented")
	} else {
		err = errors.New("credentials are not enough")
	}
	return
}

func NewConfig(creds *credentials.Credentials, region string) (config *aws.Config) {
	config = aws.NewConfig().
		WithRegion(region).
		WithCredentials(creds)
	return
}

func NewEc2Service(sess *session.Session, config *aws.Config) (ec2Service *ec2.EC2) {
	ec2Service = ec2.New(sess, config)
	return
}

func NewCloudWatchService(sess *session.Session, config *aws.Config) (cloudWatchService *cloudwatch.CloudWatch) {
	cloudWatchService = cloudwatch.New(sess, config)
	return
}
