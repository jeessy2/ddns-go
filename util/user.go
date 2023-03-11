package util

import (
	"os"
)

const ConfigFilePathENV = "DDNS_CONFIG_FILE_PATH"

// GetConfigFilePath 获得配置文件路径
func GetConfigFilePath() string {
	configFilePath := os.Getenv(ConfigFilePathENV)
	if configFilePath != "" {
		return configFilePath
	}
	return GetConfigFilePathDefault()
}

// GetConfigFilePathDefault 获得默认的配置文件路径
func GetConfigFilePathDefault() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		// log.Println("Getting Home directory failed: ", err)
		return "../.ddns_go_config.yaml"
	}
	return dir + string(os.PathSeparator) + ".ddns_go_config.yaml"
}
