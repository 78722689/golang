package analyzer

import (
	//"htmlparser"
	"routingpool"
	"utility"
)

var logger = utility.GetLogger()

func StartAnalyzer()  {
	routingpool.PutTask(AnaTask)
}

func PutData(data interface{}) {
	PutMessage(data)
}
