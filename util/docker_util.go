package util

import "os"

// DockerEnvFile Docker容器中包含的文件
const DockerEnvFile string = "/.dockerenv"

// IsRunInDocker 是否在docker中运行
func IsRunInDocker() bool {
	_, err := os.Stat(DockerEnvFile)
	return err == nil
}
