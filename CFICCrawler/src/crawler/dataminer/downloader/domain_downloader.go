package downloader

import (
	"github.com/spf13/viper"
	"htmlparser"
	"crawler/dataminer/analyzer"
)

var(
	domain_name = "SAME_GN"
)

type DOMAIN struct {

}

func (d *DOMAIN) Download(stockNumber string, stockName string, moduleURL string) {
	logger.Debugf("Domain-Downloader start to download for %s(%s}", stockName, stockNumber)
	//pageID := strings.Split(moduleURL, "/")[2]

	fileToWrite := viper.GetString("global.download_folder") + stockNumber + "/modules/" + d.ModuleName() +  "/" +d.ModuleName() + ".html"
	if err := StartDownload(viper.GetString("global.quote_homepage")+moduleURL, fileToWrite, viper.GetBool("module.jjcc.overwrite")); err != nil {
		logger.Errorf("Download Domain module failure for %s_%s.", stockName, stockNumber)
		return
	}
	logger.Debugf("Begin to parse domains for %s(%s) url:%s", stockName, stockNumber, moduleURL)
	domains := d.getDomains(fileToWrite)
	analyzer.PushDomains(map[string][]string{stockNumber:domains})
	logger.Debugf("Found domains %v", domains)
}

func (d *DOMAIN) ModuleName() string {
	return domain_name
}

func (d *DOMAIN) getDomains(file string) []string {
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		logger.Errorf("Parse file failure, %s", err)
		return nil
	}

	return doc.GetDomains()
}