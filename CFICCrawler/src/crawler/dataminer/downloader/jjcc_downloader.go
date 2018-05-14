package downloader

import (
	"utility"
	"github.com/spf13/viper"
	"htmlparser"
	"time"
	"fmt"
	"strings"
)

var(
	logger = utility.GetLogger()
	jjcc_name = "JJCC"
	)
type JJCC struct {

}

func (jjcc *JJCC) Download(stockNumber string, moduleURL string) {
	logger.Debugf("JJCC Downloader to download url=%s, file=%s", stockNumber, moduleURL)

	stockID := strings.Split(moduleURL, "/")[2]

	// 1. Get the start time and end time so that download the JJCC data during the start-end time.
	if startTime, endTime, err:= ParseDuration(viper.GetString("module.jjcc.durations")); err != nil {
		logger.Error(err)
		return
	} else {
		logger.Debugf("start=%s, end=%s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

		fileToWrite := viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() +  "/" +jjcc.ModuleName() + ".html"

		// 2. Download the first page of JJCC to parse all the record date
		if err := StartDownload(viper.GetString("global.quote_homepage")+moduleURL, fileToWrite, viper.GetBool("module.jjcc.overwrite")); err != nil {
			logger.Errorf("Download JJCC module for %s failure.", stockNumber)
			return
		}



		// 3. Find out all JJCC records date during the start-end time
		for _, recordDate := range jjcc.getAllRecordsDate(fileToWrite) {
			//d, _ := time.Parse("2006-01-02", recordDate)
			if startTime.Before(recordDate) && endTime.After(recordDate) {
				logger.Debugf("============================Found %s", recordDate)
				// 4. Download the page for this record JJCC
				StartDownload(
					viper.GetString("global.quote_homepage") + fmt.Sprintf(viper.GetString("module.jjcc.url_path"), stockID, recordDate.Format("2006-01-02")),
						viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() +  "/" +recordDate.Format("2006-01-02") + ".html",
					viper.GetBool("module.jjcc.overwrite"))
			}
		}
	}
}

func (jjcc *JJCC) getAllRecordsDate(file string) []time.Time {
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		logger.Errorf("Parse file failure, %s", err)
		return nil
	}

	return doc.JJCC_ParseRecordsDate()
}

func (jjcc *JJCC) ModuleName() string {
	return jjcc_name
}