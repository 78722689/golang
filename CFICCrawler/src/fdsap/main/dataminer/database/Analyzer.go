package database

import (
	//"htmlparser"
	"routingpool"
	"utility"
	"github.com/spf13/viper"
)

var logger = utility.GetLogger()

func StartAnalyzer()  {
	if viper.GetBool("analyser.enable") {
		routingpool.PutTask(RedisDispatcher)
	} else {
		logger.Warning("Analyser is disabled.")
	}

}

func PutData(data interface{}) {
	PushDataIntoRedis(data)
}
