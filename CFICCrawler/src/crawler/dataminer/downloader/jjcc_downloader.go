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
	stockID string
	stockNumber string
}

func (jjcc *JJCC) Download(stockNumber string, moduleURL string) {
	logger.Debugf("JJCC Downloader to download url=%s, file=%s", stockNumber, moduleURL)

	jjcc.stockID = strings.Split(moduleURL, "/")[2]
	jjcc.stockNumber = stockNumber

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

		// 3. Find out JJCC records date during the start-end time
		for _, recordDate := range jjcc.getAllRecordsDate(fileToWrite) {
			d, _ := time.Parse("2006-01-02", recordDate)
			if startTime.Before(d) && endTime.After(d) {
				logger.Debugf("============================Found %s", recordDate)
				// 4. Download the page for this record JJCC
				jjcc.download(recordDate)
			}
		}
	}
}

func (jjcc *JJCC) download(recordDate string) {
	url := viper.GetString("global.quote_homepage") + fmt.Sprintf(viper.GetString("module.jjcc.url_path"), jjcc.stockID, recordDate)
	file := viper.GetString("global.download_folder") + jjcc.stockNumber + "/modules/" + jjcc.ModuleName() +  "/" +recordDate + ".html"
	overwrite :=  viper.GetBool("module.jjcc.overwrite")
	StartDownload(url, file, overwrite)
	//logger.Infof("xxxxxxxxxxxxxxxxxxx%s", url)
}

func (jjcc *JJCC) getAllRecordsDate(file string) []string {
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