package downloader

import (
	"httpcontroller"
	"time"
	"strings"
	"fmt"
)

type Moduler interface {
	Download(stockNumber string, moduleURL string)
	ModuleName() string
}

type DownloadError string
func (str DownloadError) Error() string {
	return fmt.Sprintf("Download error %q", string(str))
}

func StartDownload(url string, file string, overwrite bool) (err error) {
	request := httpcontroller.Request{Url : url, File : file, OverWrite : overwrite}
	_, err = request.Get()
	if err != nil {
		logger.Errorf("[%]Request to url failure, %s", err)
	}

	return
}

func ParseDuration(durations string) (startTime time.Time, endTime time.Time, err error) {
	tmp := strings.Split(durations, "~")

	if durations == "" {
		err = DownloadError("Durations did not set in config file for JJCC module.")
	} else if len(tmp) == 1{
		startTime, err = time.Parse("2006-01-02", tmp[0])
		endTime = startTime
	} else if (tmp[0] == "" && tmp[1] != "") {
		startTime, err = time.Parse("2006-01-02","1990-12-19")
		endTime, err = time.Parse("2006-01-02", tmp[1])
	} else if (tmp[0] != "" && tmp[1] == "") {
		startTime, err = time.Parse("2006-01-02", tmp[0])
		endTime = time.Now()
	} else if (tmp[0] == "" && tmp[1] == "") {
		startTime, err = time.Parse("2006-01-02","1990-12-19")
		endTime = time.Now()
	} else {
		startTime, err = time.Parse("2006-01-02", tmp[0])
		endTime, err = time.Parse("2006-01-02", tmp[1])
	}

	return
}