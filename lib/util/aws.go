package util

import (
	"../aws"
	"../reporter"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSClient struct {
	Region string
	EC2Service aws.EC2Client
	EBSService aws.EBSClient
}

func (A AWSClient) printSummary(res Resource, rep reporter.Reporter)  {
	formatting := `
============================================================
Resource Configuration Summary
============================================================
name       : %s
============================================================`
	summary := fmt.Sprintf(formatting, res.Name)
	fmt.Println(summary)
}

func (A AWSClient) clean(res Resource, rep reporter.Reporter)  {
	// ebs
	if res.Name == "ebs" {
		A.EBSService.Clean(rep)
	}


}


// Get AWS session
func _get_aws_session() *session.Session {
	mySession := session.Must(session.NewSession())
	return mySession
}

//Bootstrap process for getting right clients
func bootstrapServices(region string, assume_role string) AWSClient {
	aws_session := _get_aws_session()

	var creds *credentials.Credentials
	if len(assume_role) != 0  {
		creds = stscreds.NewCredentials(aws_session, assume_role)
	}

	//Get all clients
	client := AWSClient{
		Region: region,
		EC2Service: aws.NewEC2Client(aws_session, region, creds),
		EBSService: aws.NewEBSClient(aws_session, region, creds),
	}

	return client
}

