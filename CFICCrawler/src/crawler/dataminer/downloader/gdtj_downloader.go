package downloader

import (
	"github.com/spf13/viper"
)

//var logger = utility.GetLogger()
var gdtj_name  = "GDTJ"

type GDTJ struct {
}

func (gdtj *GDTJ) Download(stockNumber string, moduleURL string) {
	logger.Debugf("GDTJ Downloader to download  url=%s, file=%s", stockNumber, moduleURL)
	/*request := httpcontroller.Request{
		Url:        viper.GetString("global.quote_homepage") + moduleURL,
		File:       viper.GetString("global.download_folder") + stockNumber + "/modules/" + gdtj.ModuleName() + ".html",
		OverWrite:  viper.GetBool("module.gdtj.overwrite")}
	_, err := request.Get()
	if err != nil {
		logger.Errorf("[%]Request to url failure, %s", err)
	}
	*/
	StartDownload(viper.GetString("global.quote_homepage") + moduleURL,
		viper.GetString("global.download_folder") + stockNumber + "/modules/" + gdtj.ModuleName() + ".html",
			viper.GetBool("module.gdtj.overwrite"))
}

func (g *GDTJ) ModuleName() string {
	return gdtj_name
}