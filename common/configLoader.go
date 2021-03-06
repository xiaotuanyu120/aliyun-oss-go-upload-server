package common

import (
	"fmt"

	"github.com/spf13/viper"
)

var viperConfig *viper.Viper

func Init() {
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viperConfig = viper.GetViper()
}

func GetConfig() *viper.Viper {
	return viperConfig
}
