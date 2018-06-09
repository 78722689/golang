package analyzer

import (
	"routingpool"
	"htmlparser"
	"github.com/spf13/viper"
	"time"
)

var (
	AnaTask *AnalysisTask
)

func init() {
	AnaTask = NewAnalysisTask()
}

// Implement Task interface
type AnalysisTask struct {
	*routingpool.Base
	message chan interface{}
	//timer int	// The waiting seconds for receiving data, analysis routine exits after the waiting.
}

func NewAnalysisTask() *AnalysisTask {
	return &AnalysisTask{message : make(chan interface{}, 1024), Base:&routingpool.Base{Name: "Analysis Task", Response: make(chan bool)}}
}

func PutMessage(msg interface{}) {
	AnaTask.message <- msg
}

func (a *AnalysisTask) caller(id int) {
	timeout := time.NewTimer(time.Second * time.Duration(viper.GetInt("analysis.timer")))
	exit := true

	for exit {
		select {
			case data := <-a.message:
				logger.Info("analysis task received data")
				tmp := data.([]*htmlparser.JJCCData)

				for _, value := range tmp {
					logger.Infof("Row data name %s, code %s, holdcount %.4f, holdvalue %.4f", value.Name, value.Code, value.HoldCount, value.HoldValue)
				}

				timeout.Reset(time.Second * time.Duration(viper.GetInt("analysis.timer")))

			case <- timeout.C: // The waiting seconds for receiving data, analysis routine exits after the waiting.
				logger.Info("Analyser routine exit.....................")
				exit = false
				break

		}
	}
}

func (a *AnalysisTask) Run(id int) {
	a.caller(id)
}