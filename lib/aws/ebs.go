package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ebs"
	"github.com/aws/aws-sdk-go/service/ec2"
	Logger "github.com/sirupsen/logrus"
	"strings"

	"../reporter"
)

type EBSClient struct {
	Client *ebs.EBS
	EC2Client *ec2.EC2
}

type Wasted struct {
	Id string
	Ctime string
}

type TestV struct {
	Id string
	Key string
	Value string
}

// Clean wasted EBS Volumes
func (e EBSClient) Clean(rep reporter.Reporter) {

	_logging_with_slack(rep, "Start cleaning `wasted EBS volumes`")

	// Find all volumes
	volumes, err := e.EC2Client.DescribeVolumes(&ec2.DescribeVolumesInput{})
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
		return
	}

	wasted := []Wasted{}
	testv := []TestV{}
	avail_count := 0
	test_count := 0
	for _, v := range volumes.Volumes {
		if *v.State == "available" {
			avail_count += 1
			wasted = append(wasted, Wasted{
				Id:    *v.VolumeId,
				Ctime: fmt.Sprint(v.CreateTime),
			})
			continue
		}

		for _, t := range v.Tags {
			if strings.Contains(strings.ToLower(*t.Value), "test") {
				test_count += 1
				testv = append(testv, TestV{
					Id:    *v.VolumeId,
					Key:   *t.Key,
					Value: *t.Value,
				})
				break
			}
		}
	}

	_print_wasted_volumes(avail_count, wasted)
	_print_test_volumes(test_count, testv)


	rep.Send_slack_message(fmt.Sprintf("[EBS volumes] Wasted: %d / Test Tagged: %d", avail_count, test_count))
}

func NewEBSClient(session *session.Session, region string, creds *credentials.Credentials) EBSClient {
	return EBSClient{
		Client: _get_ebs_client_fn(session, region, creds),
		EC2Client: _get_ec2_client_fn(session, region, creds),
	}
}

func _get_ebs_client_fn(session *session.Session, region string, creds *credentials.Credentials) *ebs.EBS {
	if creds == nil {
		return ebs.New(session, &aws.Config{Region: aws.String(region)})
	}
	return ebs.New(session, &aws.Config{Region: aws.String(region), Credentials: creds})
}

func _print_wasted_volumes(count int, volumes []Wasted)  {
	Logger.WithFields(Logger.Fields{
		"count": count,
	}).Info("Cleaner found unattached volumes")

	if count <= 0 { return }
	for i, v := range volumes {
		Logger.WithFields(Logger.Fields{
			"Volume ID": v.Id,
			"Created At": v.Ctime,
		}).Info("Volume ", i+1)
	}
}

func _print_test_volumes(count int, volumes []TestV)  {
	Logger.WithFields(Logger.Fields{
		"count": count,
	}).Info("Cleaner found volumes with test")

	if count <= 0 { return }
	for i, v := range volumes {
		Logger.WithFields(Logger.Fields{
			"Volume ID": v.Id,
			"Key": v.Key,
			"Value": v.Value,
		}).Info("Volume ", i+1)
	}
}

func _logging_with_slack(rep reporter.Reporter, msg string)  {
	Logger.Info(msg)
	rep.Send_slack_message(msg)
}
