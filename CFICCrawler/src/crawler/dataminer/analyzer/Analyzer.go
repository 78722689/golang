package analyzer

import (
	//"htmlparser"
	"routingpool"
	"utility"
	"htmlparser"
)

type Analyer struct {
	data         chan interface{}
}

var logger = utility.GetLogger()
var a *Analyer

func init()  {
	a = new(Analyer)
	a.data = make(chan interface{})
}

func StartAnalyzer()  {
	routingpool.PutTask(routingpool.NewCaller("Module analyzer", a.caller))
}

func PutData(data interface{}) {
	a.data <- data
}

func (r *Analyer) caller(id int) {
	for {
		data := <-r.data
		tmp := data.([]*htmlparser.JJCCData)

		for _, value := range tmp {
			logger.Infof("Received data name %s, code %s, holdcount %.4f, holdvalue %.4f", value.Name, value.Code, value.HoldCount, value.HoldValue)
		}
	}
}