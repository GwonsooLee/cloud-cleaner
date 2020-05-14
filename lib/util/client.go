package util


import (
	"../reporter"
	Logger "github.com/sirupsen/logrus"
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
	reporter := reporter.New(slackConfig.Token, slackConfig.ChannelId)

	// Print summary of resource to clean
	c.printSummary(resource, reporter)

	//start clean
	c.clean(resource, reporter)
}