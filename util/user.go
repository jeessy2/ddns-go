package util

import (
	"log"
	"os"
	"os/user"
)

// GetConfigFilePath 获得配置文件路径
func GetConfigFilePath() string {
	u, err := user.Current()
	if err != nil {
		log.Println("Geting current user failed!")
		return "../.ddns_go_config.yaml"
	}
	return u.HomeDir + string(os.PathSeparator) + ".ddns_go_config.yaml"
}
