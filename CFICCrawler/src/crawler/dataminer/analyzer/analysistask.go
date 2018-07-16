package analyzer

import (
	"routingpool"
	"htmlparser"
	"github.com/spf13/viper"
	"time"
	"os"
	"bufio"
	"github.com/axgle/mahonia"
	"strings"
	"utility"
	"fmt"
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
	timeout := time.NewTimer(time.Second * time.Duration(viper.GetInt("analyser.timer")))
	exit := true
	funds, _ := a.getFunds()

	for exit {
		select {
			case data := <-a.message:
				logger.Info("analysis task received data")
				tmp := data.([]*htmlparser.JJCCData)

				for _, value := range tmp {
					logger.Infof("Row data name %s, code %s, holdcount %.4f, holdvalue %.4f", value.Name, value.Code, value.HoldCount, value.HoldValue)
					for _,fund := range funds {
						if strings.Contains(value.Name, fund) {
							filename := fmt.Sprintf("%s%s.csv",viper.GetString("global.output_folder"),value.Name)
							line := fmt.Sprintf("%s, %s, %.4f, %.4f", value.Code, value.RecordDate, value.HoldCount, value.HoldValue)
							utility.WriteToFile(filename, line)
						}
					}
				}

				timeout.Reset(time.Second * time.Duration(viper.GetInt("analyser.timer")))

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

func (a *AnalysisTask) getFunds() ([]string, error) {
	filename := viper.GetString("module.jjcc.funds_file_path")
	file, err:= os.Open(filename)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "\n", filename, err)
		logger.Errorf("WARN: Open file %s failed, %s", filename, err)
		return nil, err
	}
	defer file.Close()

	var result []string
	decoder := mahonia.NewDecoder("gbk")
	scanner := bufio.NewScanner(decoder.NewReader(file))
	for scanner.Scan() {
		//fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
		result = append(result, scanner.Text())
	}

	return result,nil
}