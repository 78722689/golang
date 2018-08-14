package dataminer

import (
	"htmlparser"
	"httpcontroller"
	"os"
	"strings"
	"utility"
	"crawler/dataminer/downloader"
	"github.com/spf13/viper"
	"crawler/dataminer/database"
)

var logger = utility.GetLogger()


func (t *Target) RegisterModuleDownloader(m downloader.Moduler) *Target{
	t.Modules = append(t.Modules, m)

	return t
}

func (t *Target) Start() {

	// Request stock list page to get all the stocks
	request := httpcontroller.Request{
		Url:    viper.GetString("global.quote_homepage") + viper.GetString("global.stock_list_url_path"),
	}
	root, _ := request.Get()

	doc, err := htmlparser.ParseFromNode(root)
	if err != nil {
		logger.Errorf("Parse file error, %v", err)

		os.Exit(1)
	}

	//mainTask := downloader.NewDownloadTask("Main Downloader", func(id int) {
	func(id int) {
		for _, stockinfo := range doc.GetStocks(t.Stocks) {
			// Push stock name to table
			database.PushStocks([]string{"STOCKS_NAME_SHH", stockinfo.Name + "_" + stockinfo.Number})

			// To request home page.
			tempStockinfo := stockinfo // Copy the value, so that below closure run correctly
			//homepageCaller := func(id int) {
			func(id int) {
				logger.Infof("[Thread-%d] Downloading link:%v name:%v, number:%v", id, tempStockinfo.Link, tempStockinfo.Name, tempStockinfo.Number)

				file := viper.GetString("global.download_folder") + tempStockinfo.Number + "/" + tempStockinfo.Link
				homepage_request := httpcontroller.Request{
					Url:       viper.GetString("global.quote_homepage") + "/" + tempStockinfo.Link,
					File:      file,
					OverWrite: true}
				_, err := homepage_request.Get()
				if err != nil {
					logger.Errorf("Request failure, %s", err)
					return
				}

				logger.Infof("[Thread-%d] Downloaded homepage link:%v name:%v, number:%v", id, tempStockinfo.Link, tempStockinfo.Name, tempStockinfo.Number)

				// To request modules for each stock.
				//moduleCaller := func(id int) {
				func(id int) {
					doc, err := htmlparser.ParseFromFile(file)
					if err != nil {
						logger.Errorf("Parse file failure, %s", err)
						return
					}

					stock_modules_url := doc.GetModuleURL(tempStockinfo.Link)
					for _, url := range stock_modules_url {
						for _, module := range t.Modules {
							if strings.Contains(strings.ToUpper(url), "/"+module.ModuleName()+"/") {
								module.Download(tempStockinfo.Number, tempStockinfo.Name, url)
								break
							}
						}
					}
				}(777)
				//routingpool.PutTask(downloader.NewDownloadTask("Modules Downloader", moduleCaller))
			}(888)
			//routingpool.PutTask(downloader.NewDownloadTask("Homepage Downloader", homepageCaller))
		}
	}(999)
		//})
	//routingpool.PutTask(mainTask)
}
