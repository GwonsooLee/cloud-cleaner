package lib

import (
	"./util"
	"./reporter"
)


type Runner struct {
	Config util.Config
	CheckList util.CheckList
}

func Run()  {
	// parse argument from command
	config := util.ParseArgument()

	//Read Checklist to clean from manifest file
	checkList := util.ReadConfigFromManifest(config.ConfigPath)
	checkList.Validate()

	//Get Runner Struct
	runner := Runner{
		Config:    config,
		CheckList: checkList,
	}

	runner.Clean()
}

func (r Runner) Clean()  {
	resources := r.CheckList.Resources
	regions := r.CheckList.Regions

	reporter := reporter.New(r.CheckList.SlackConfig.Token, r.CheckList.SlackConfig.ChannelId)
	reporter.SendTitleMessage()

	for _, region := range regions {
		awsClient := util.NewAWSClient(region, r.Config.AssumeRole)
		reporter.SendRegionMessage(region)
		for _, resource := range resources {
			util.Start(awsClient, region, resource, reporter)
		}
	}

}

