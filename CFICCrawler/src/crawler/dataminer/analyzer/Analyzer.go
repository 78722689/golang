package analyzer

import (
	//"htmlparser"
	"routingpool"
	"utility"
	"github.com/spf13/viper"
)

var logger = utility.GetLogger()

func StartAnalyzer()  {
	if viper.GetBool("analyser.enable") {
		routingpool.PutTask(AnaTask)
	} else {
		logger.Warning("Analyser is disabled.")
	}

}

func PutData(data interface{}) {
	PutMessage(data)
}
