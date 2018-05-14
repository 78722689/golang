package utility

import (
	"github.com/spf13/viper"
	"fmt"
)

func NewConfig(configFilePath string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//fmt.Println(viper.GetBool("module.gdtj.overwrite"))
}