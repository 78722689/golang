package downloader

import (
	"utility"
	"github.com/spf13/viper"
	"htmlparser"
	"time"
	"fmt"
	"strings"
	"os"
	"crawler/dataminer/collector"
)

var(
	logger = utility.GetLogger()
	jjcc_name = "JJCC"
	)
type JJCC struct {

}

func (jjcc *JJCC) Download(stockNumber string, stockName string, moduleURL string) {
	pageID := strings.Split(moduleURL, "/")[2]

	logger.Debugf("JJCC-Downloader start to download for stockNumber=%s, pageID=%s, moduleUrl=%s", stockNumber, pageID, moduleURL)

	// 1. Get the start time and end time so that download the JJCC data during the start-end time.
	if startTime, endTime, err:= ParseDuration(viper.GetString("module.jjcc.durations")); err != nil {
		logger.Error(err)
		return
	} else {
		//logger.Debugf("start=%s, end=%s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

		fileToWrite := viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() +  "/" +jjcc.ModuleName() + ".html"
		// 2. Download the first page of JJCC to parse all the record date
		if err := StartDownload(viper.GetString("global.quote_homepage")+moduleURL, fileToWrite, viper.GetBool("module.jjcc.overwrite")); err != nil {
			logger.Errorf("Download JJCC module for %s failure.", stockNumber)
			return
		}

		// 3. Find out all JJCC records date during the start-end time
		for index, recordDate := range jjcc.getAllRecordsDate(fileToWrite) {
			if startTime.Before(recordDate) && endTime.After(recordDate) {
				logger.Debugf("Matched JJCC record on date %s for stockNumber %s(%s) ", recordDate.Format("2006-01-02"), stockNumber, pageID)

				file := viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() + "/" + recordDate.Format("2006-01-02") + ".html"

				// Sometimes the newest JJCC record can not be opened on Web Browser due to bug in QUOTE
				// And the newest record has been download before, so here do not need to download it again.
				// So rename the page is enough.
				if index == 0 {
					//newName := viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() +  "/" + recordDate.Format("2006-01-02") + ".html"
					os.Rename(fileToWrite, file)
				} else {
					//file := viper.GetString("global.download_folder") + stockNumber + "/modules/" + jjcc.ModuleName() + "/" + recordDate.Format("2006-01-02") + ".html"
					// 4. Download the page for this JJCC record
					StartDownload(
						viper.GetString("global.quote_homepage")+fmt.Sprintf(viper.GetString("module.jjcc.url_path"), pageID, recordDate.Format("2006-01-02")),
						file,
						viper.GetBool("module.jjcc.overwrite"))
				}

				// If Analyser is not enabled, only download the pages
				if viper.GetBool("analyser.enable") {
					collector.CollectJJCC(stockNumber, stockName, file, recordDate.Format("2006-01-02"))
				}
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