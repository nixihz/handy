package utils

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

//初始化配置文件
func NewConfig(fileName string) {

	filePath := GetRunPath()

	Config = viper.New()
	Config.SetConfigType("toml")
	Config.AutomaticEnv()
	Config.SetConfigName(fileName)
	Config.AddConfigPath(filePath)
	Config.AddConfigPath("configs")
	Config.AddConfigPath("build/configs")
	//Config.WatchConfig()

	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
}
