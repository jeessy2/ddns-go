package util

import (
	"log"
	"os"
	"os/user"
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
	u, err := user.Current()
	if err != nil {
		log.Println("Geting current user failed!")
		return "../.ddns_go_config.yaml"
	}
	return u.HomeDir + string(os.PathSeparator) + ".ddns_go_config.yaml"
}
