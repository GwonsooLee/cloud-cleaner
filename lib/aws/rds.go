package aws

import (
	"../reporter"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	Logger "github.com/sirupsen/logrus"
	"strings"
)

type RDSClient struct {
	Client *rds.RDS
}

type RDSPrintable struct {
	DBInstanceIdentifier  string
	DBInstanceClass  	  string
	InstanceCreateTime	  string
}

func (r RDSClient) Clean(rep reporter.Reporter) {
	// Detect Wasted Volumes
	Logger.Info("Start cleaning `test tagged`")

	testTagged := _findTestTaggedInstance(r.Client)
	printable := _makePrintable(testTagged)
	_print_test_tagged_instances(len(printable), printable, rep)
}


func NewRDSClient(session *session.Session, region string, creds *credentials.Credentials) RDSClient {
	return RDSClient{
		Client: _get_rds_client_fn(session, region, creds),
	}
}

func _get_rds_client_fn(session *session.Session, region string, creds *credentials.Credentials) *rds.RDS {
	if creds == nil {
		return rds.New(session, &aws.Config{Region: aws.String(region)})
	}
	return rds.New(session, &aws.Config{Region: aws.String(region), Credentials: creds})
}

func _findTestTaggedInstance(c *rds.RDS) []*rds.DBInstance {
	input := &rds.DescribeDBInstancesInput{}

	output, err := c.DescribeDBInstances(input)
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

	filtered := []*rds.DBInstance{}
	for _, instance := range output.DBInstances {
		fmt.Println(instance)
		if strings.Contains(*instance.DBInstanceIdentifier, "test") {
			filtered = append(filtered, instance)
		}
	}

	return filtered
}


func _makePrintable(raw []*(rds.DBInstance)) []RDSPrintable {
	printable := []RDSPrintable{}
	for _, instance := range raw {
		printable = append(printable, RDSPrintable{
			DBInstanceIdentifier:   *instance.DBInstanceIdentifier,
			DBInstanceClass:    	*instance.DBInstanceClass,
			InstanceCreateTime: 	fmt.Sprintf("%v",*instance.InstanceCreateTime),
		})
	}

	return printable
}

func _print_test_tagged_instances(count int, instances []RDSPrintable, rep reporter.Reporter)  {
	Logger.WithFields(Logger.Fields{
		"count": count,
	}).Info("Cleaner found test tagged instances")

	if count <= 0 { return }

	textList := []string{}
	for idx, instance := range instances {
		Logger.WithFields(Logger.Fields{
			"Name": instance.DBInstanceIdentifier,
			"Class": instance.DBInstanceClass,
			"Created At": instance.InstanceCreateTime,
		}).Info("Test Tagged ", idx+1)

		textList = append(textList, fmt.Sprintf("DB-Name=`%s` Class=`%s` Created At=%s", instance.DBInstanceIdentifier, instance.DBInstanceClass, instance.InstanceCreateTime))
	}

	title := CRITICAL+"*Please check these `test tagged RDS instance`!!*"
	msgOption := rep.CreateAlarmMessage(title, textList)
	rep.SendBlockMessage(msgOption)
}
