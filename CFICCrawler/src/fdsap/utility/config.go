package utility

import (
	"github.com/spf13/viper"
    //"github.com/op/go-logging"
	"fmt"
    //"fdsap/http"
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