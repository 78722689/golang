package dataminer

import (
	"htmlparser"
	"httpcontroller"
	"os"
	"routingpool"
	"strings"
	"utility"
	"crawler/dataminer/downloader"
	"github.com/spf13/viper"
)

var logger = utility.GetLogger()

/*
type DownloadInfo struct {
	Folder      string
	Proxy       *httpcontroller.Proxy
	Overwrite   bool
	Stocks      []string // Empty for downloading all stocks.
	CodeChannel chan string

	RoutingPool *routingpool.ThreadPool
}
*/
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

	mainTask := routingpool.NewCaller("Main Downloader", func(id int) {
		for _, stockinfo := range doc.GetStocks(t.Stocks) {

			// To request home page.
			tempStockinfo := stockinfo // Copy the value, so that below closure run correctly
			homepageCaller := func(id int) {
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
				moduleCaller := func(id int) {
					doc, err := htmlparser.ParseFromFile(file)
					if err != nil {
						logger.Errorf("Parse file failure, %s", err)
						return
					}

					stock_modules_url := doc.GetModuleURL(tempStockinfo.Link)
					for _, url := range stock_modules_url {
						for _, module := range t.Modules {
							if strings.Contains(strings.ToUpper(url), "/" + module.ModuleName() + "/") {
								module.Download(tempStockinfo.Number, url)
								break
							}
						}
					}

					// Download history data.
					//histData := modulehandler.HTD{Code: tempStockinfo.Number, Folder: t.Folder, Proxy: t.Proxy}
					//if err := histData.Download(); err != nil {
					//	logger.ERROR(fmt.Sprintf("Download history data failure, %s", err))
					//}

					//syncChan <- true
				}

				// Start the Miner to collect/analyze data
				/*collector := Collect{
					Code:        tempStockinfo.Number,
					Folder:      t.Folder,
					SyncChan:    syncChan,
					RoutingPool: t.RoutingPool,
					Proxy:       t.Proxy}
				collector.Start()
				*/
				//t.StartAnalyse(tempStockinfo.Number, syncChan)

				routingpool.PutTask(routingpool.NewCaller("Modules Downloader", moduleCaller))
			}

			routingpool.PutTask(routingpool.NewCaller("Homepage Downloader", homepageCaller))
		}

	})

	routingpool.PutTask(mainTask)
}
