package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Client struct {
	Client *ec2.EC2
	AsClient *autoscaling.AutoScaling
}

func NewEC2Client(session *session.Session, region string, creds *credentials.Credentials) EC2Client {
	return EC2Client{
		Client: _get_ec2_client_fn(session, region, creds),
		AsClient: _get_asg_client_fn(session, region, creds),
	}
}

func _get_ec2_client_fn(session *session.Session, region string, creds *credentials.Credentials) *ec2.EC2 {
	if creds == nil {
		return ec2.New(session, &aws.Config{Region: aws.String(region)})
	}
	return ec2.New(session, &aws.Config{Region: aws.String(region), Credentials: creds})
}


func _get_asg_client_fn(session *session.Session, region string, creds *credentials.Credentials) *autoscaling.AutoScaling {
	if creds == nil {
		return autoscaling.New(session, &aws.Config{Region: aws.String(region)})
	}
	return autoscaling.New(session, &aws.Config{Region: aws.String(region), Credentials: creds})
}

