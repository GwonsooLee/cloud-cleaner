package util

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	Logger "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	ConfigPath string
	AssumeRole string
}

type CheckList struct {
	IgnoreTag string `yaml:"ignore_tag"`
	Regions []string `yaml:"regions"`
	SlackConfig SlackConfig `yaml:"slack"`
	Resources []Resource
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Token string `yaml:"token"`
}

type Resource struct {
	Name string
}

func (c CheckList) Validate()  {

	// At least one resource should be specified
	if len(c.Resources) <= 0{
		Logger.Error(NO_RESOURCE_SPECIFIED)
		os.Exit(1)
	}

	// At least one region should be specified
	if len(c.Regions) <= 0{
		Logger.Error(NO_REGION_SPECIFIED)
		os.Exit(1)
	}

	// Add more for validation...
}

// Parsing Config from command
func ParseArgument() Config {
	// Parsing Argument
	Logger.Info("Checking out arguments for command...")
	//region := flag.String("region", "ap-northeast-2", "AWS Region ID to check resources and clean")
	configPath := flag.String("config", "checklist.yaml", "Path of cleaner config yaml file")
	assumeRole := flag.String("assume-role", "", "ARN of assume role")

	flag.Parse()

	config := Config{
		//Region: *region,
		ConfigPath: *configPath,
		AssumeRole: *assumeRole,
	}

	return config
}


// Parsing Manifest File
func ReadConfigFromManifest(manifest string) CheckList {

	Logger.Info("Reading Configuration File...")
	checkList := CheckList{}

	//Check File exists
	if ! fileExists(manifest) {
		Logger.Error(CONFIG_FILE_NOT_EXIST)
	}

	yamlFile, err := ioutil.ReadFile(manifest)
	if err != nil {
		Logger.Error("Failed to read yaml file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &checkList)
	if err != nil {
		Logger.Error("Failed to unmarshal data: %v", err)
	}

	return checkList
}
