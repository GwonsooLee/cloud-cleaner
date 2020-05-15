package aws

import (
	"../reporter"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	Logger "github.com/sirupsen/logrus"
)

type EC2Client struct {
	Client *ec2.EC2
	AsClient *autoscaling.AutoScaling
}

type Filtered struct {
	InstanceId 	  string
	InstanceType  string
	LaunchTime	  string
}

func (e EC2Client) Clean(rep reporter.Reporter) {
	// Detect Wasted Volumes
	Logger.Info("Start cleaning `untagged and stopped EC2 volumes`")

	stopped := findStoppedInstance(e.Client)
	filtered := _get_filtered_list(stopped)
	_print_stopped_instances(len(filtered), filtered, rep)

	untagged := findUntaggedInstance(e.Client, []*ec2.Instance{}, nil)
	filtered = _get_filtered_list(untagged)
	_print_untagged_instances(len(filtered), filtered, rep)

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

func findStoppedInstance(c *ec2.EC2) []*(ec2.Instance) {
	// Find stopped instances
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-code"),
				Values: []*string{
					aws.String("80"),
				},
			},
		},
	}

	instances, err := c.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	reserveGroup := []*ec2.Instance{}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			reserveGroup = append(reserveGroup, instance)
		}
	}

	return reserveGroup
}

func findUntaggedInstance(c *ec2.EC2, reserveGroup []*(ec2.Instance), nextToken *string) []*(ec2.Instance) {

	// Find stopped instances
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-code"),
				Values: []*string{
					aws.String("16"),
				},
			},
		},
		NextToken: nextToken,
	}

	instances, err := c.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			untagged := true
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" && len(*tag.Value) >0 {
					untagged = false
					break
				}
			}

			if untagged { reserveGroup = append(reserveGroup, instance) }
		}
	}


	if instances.NextToken != nil {
		return findUntaggedInstance(c, reserveGroup, instances.NextToken)
	}

	return reserveGroup
}


func _print_stopped_instances(count int, instances []Filtered, rep reporter.Reporter)  {
	Logger.WithFields(Logger.Fields{
		"count": count,
	}).Info("Cleaner found stopped instances")

	if count <= 0 { return }

	textList := []string{}
	for idx, instance := range instances {
		Logger.WithFields(Logger.Fields{
			"Instance ID": instance.InstanceId,
			"Instance Type": instance.InstanceType,
			"Launched At": instance.LaunchTime,
		}).Info("Stopped ", idx+1)

		textList = append(textList, fmt.Sprintf("Instance-ID=`%s` Type=`%s` Launched At=%s", instance.InstanceId, instance.InstanceType, instance.LaunchTime))
	}

	title := CRITICAL+"*Please check these `stopped instances`!!*"
	msgOption := rep.CreateAlarmMessage(title, textList)
	rep.SendBlockMessage(msgOption)
}


func _print_untagged_instances(count int, instances []Filtered, rep reporter.Reporter)  {
	Logger.WithFields(Logger.Fields{
		"count": count,
	}).Info("Cleaner found untagged instances")

	if count <= 0 { return }

	textList := []string{}
	for idx, instance := range instances {
		Logger.WithFields(Logger.Fields{
			"Instance ID": instance.InstanceId,
			"Instance Type": instance.InstanceType,
			"Launched At": instance.LaunchTime,
		}).Info("Untagged ", idx+1)

		textList = append(textList, fmt.Sprintf("Instance-ID=`%s` Type=`%s` Launched At=%s", instance.InstanceId, instance.InstanceType, instance.LaunchTime))
	}

	title := CRITICAL+"*Please check these `untagged instances`!!*"
	msgOption := rep.CreateAlarmMessage(title, textList)
	rep.SendBlockMessage(msgOption)
}

func _get_filtered_list(raw []*(ec2.Instance)) []Filtered {
	filtered := []Filtered{}
	for _, instance := range raw {
		filtered = append(filtered, Filtered{
			InstanceId:   *instance.InstanceId,
			InstanceType: *instance.InstanceType,
			LaunchTime:   fmt.Sprintf("%v", *instance.LaunchTime),
		})
	}

	return filtered
}
