package downloader

import (
	"dataminer"
	"fmt"
	"htmlparser"
	"httpcontroller"
	"os"
	"routingpool"
	"strings"
	"time"
	"utility"
)

var logger = utility.GetLogger()

const (
	STOCK_LIST_URL     string = "http://quote.cfi.cn/stockList.aspx?t=11"
	QUOTE_HOMEPAGE_URL string = "http://quote.cfi.cn/"
)

type DownloadInfo struct {
	Folder    string
	Proxy     *httpcontroller.Proxy
	Overwrite bool

	RoutingPool *routingpool.ThreadPool
}

type DownloadTask struct {
	Name string
}

func (info *DownloadInfo) Task(id int) {
	//fmt.Println(fmt.Sprintf("Thread - %d is running with DownloadTask - %s", id, task.Name))
	time.Sleep(time.Second * 2)
}

func (d *DownloadInfo) DownloadAll() {
	d.downloadStocksHomePage([]string{}, d.Overwrite)
	//d.downloadModules([]string{}, false)
}

func (d *DownloadInfo) DownloadByStockIDs(ids []string) {
	d.downloadStocksHomePage(ids, d.Overwrite)
	//d.downloadModules([]string{}, true)
}

// Download the stocks home page, if overwrite is true, the exist home page will be rewrite.
func (d *DownloadInfo) downloadStocksHomePage(ids []string, overwrite bool) {
	// Request homepage to get all the stocks
	request := httpcontroller.Request{
		Url:   STOCK_LIST_URL,
		Proxy: d.Proxy,
	}
	root, _ := request.Get()

	doc, err := htmlparser.ParseFromNode(root)
	if err != nil {
		logger.ERROR(fmt.Sprintf("Parse file error, %v", err))

		os.Exit(1)
	}

	mainTask := routingpool.NewCaller("Main Downloader", func(id int) {
		for _, stockinfo := range doc.GetStocks(ids) {

			// To request home page.
			tempStockinfo := stockinfo // Copy the value, so that below closure run correctly
			homepageCaller := func(id int) {
				logger.INFO(fmt.Sprintf("[Thread-%d] Downloading link:%v name:%v, number:%v", id, tempStockinfo.Link, tempStockinfo.Name, tempStockinfo.Number))

				file := d.Folder + tempStockinfo.Number + "/" + tempStockinfo.Link
				stock_request := httpcontroller.Request{
					Url:       QUOTE_HOMEPAGE_URL + tempStockinfo.Link,
					File:      file,
					Proxy:     d.Proxy,
					OverWrite: overwrite}
				_, err := stock_request.Get()
				if err != nil {
					logger.ERROR(fmt.Sprintf("Request failure, %s", err))
					return
				}

				logger.INFO(fmt.Sprintf("[Thread-%d] Downloaded homepage link:%v name:%v, number:%v", id, tempStockinfo.Link, tempStockinfo.Name, tempStockinfo.Number))

				syncChan := make(chan bool)
				// To request modules for each stock.
				moduleCaller := func(id int) {
					doc, err := htmlparser.ParseFromFile(file)
					if err != nil {
						logger.ERROR(fmt.Sprintf("Parse file failure, %s", err))
						return
					}

					stock_modules_url := doc.GetModuleURL(tempStockinfo.Link)
					for _, url := range stock_modules_url {
						// an example filter on a stock modules
						filters := []string{"gdtj", "fhpx"}

						var found bool = false
						for _, filter := range filters {
							if strings.Contains(url, filter) {
								found = true
								break
							}
						}

						if found {
							logger.INFO(fmt.Sprintf("[Thread-%d] Found module, requesting %s\n", id, url))

							values := strings.Split(url, "/")
							var moduleName string
							if len(values) == 4 {
								moduleName = values[1] // eg. "/tzzk/19770/601015.html"
							} else {
								moduleName = values[2] // eg. "http://gg.cfi.cn/cbgg/19770/601015.html"
							}

							file := d.Folder + tempStockinfo.Number + "/modules/" + moduleName + ".html"
							request := httpcontroller.Request{
								Proxy:     d.Proxy,
								Url:       QUOTE_HOMEPAGE_URL + url,
								File:      file,
								OverWrite: overwrite}
							_, err := request.Get()
							if err != nil {
								logger.ERROR(fmt.Sprintf("Request to url failure, %s", err))
								continue
							}
						}
					}

					syncChan <- true
				}

				// Start the Miner to collect/analyze data
				miner := dataminer.Target{Code: tempStockinfo.Number,
					Folder:      d.Folder,
					SyncChan:    syncChan,
					RoutingPool: d.RoutingPool,
					Proxy:       d.Proxy}
				miner.Start()

				d.RoutingPool.PutTask(routingpool.NewCaller("Modules Downloader", moduleCaller))
			}

			d.RoutingPool.PutTask(routingpool.NewCaller("Stock Homepage Downloader", homepageCaller))
		}
	})

	d.RoutingPool.PutTask(mainTask)
}
