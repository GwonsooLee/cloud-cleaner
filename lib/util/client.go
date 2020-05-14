package util


import (
	"../reporter"
	Logger "github.com/sirupsen/logrus"
)

type Client interface {
	printSummary(res Resource, rep reporter.Reporter, region string)
	clean(res Resource, rep reporter.Reporter, region string)
}

func NewAWSClient(region, assume_role string) AWSClient {
	Logger.Info("Building new AWS client...")
	return bootstrapServices(region, assume_role)
}

// Start Resource cleaning
func Start(c Client, region string, resource Resource, reporter reporter.Reporter) {

	// Print summary of resource to clean
	c.printSummary(resource, reporter, region)

	//start clean
	c.clean(resource, reporter, region)
}