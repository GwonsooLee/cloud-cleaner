package util

import "os"

var (
	CONFIG_FILE_NOT_EXIST="Configuration File does not exists"
	NO_RESOURCE_SPECIFIED="No Resource is specified in checklist YAML file"
	NO_REGION_SPECIFIED="No Region is specified in checklist YAML file"
	)


// Check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

