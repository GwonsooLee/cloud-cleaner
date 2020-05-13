package util


import (
	Logger "github.com/sirupsen/logrus"
	"../reporter"
)

type Client interface {
	printSummary(res Resource, rep reporter.Reporter)
	clean(res Resource, rep reporter.Reporter)
}

func NewAWSClient(region, assume_role string) AWSClient {
	Logger.Info("Building new AWS client...")
	return bootstrapServices(region, assume_role)
}

// Start Resource cleaning
func Start(c Client, resource Resource, slackConfig SlackConfig) {
	// Make new reporter
	reporter := reporter.New(slackConfig.WebhookURL, slackConfig.Token)

	// Print summary of resource to clean
	c.printSummary(resource, reporter)

	//start clean
	c.clean(resource, reporter)
}