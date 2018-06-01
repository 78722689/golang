package crawler

import (
	"crawler/dataminer"
	"crawler/dataminer/downloader"
	"crawler/dataminer/analyzer"
)

func StartCrawl(stocks []string) {

	/*
		collector := Collect{
			Code:        tempStockinfo.Number,
			Folder:      t.Folder,
			SyncChan:    syncChan,
			RoutingPool: t.RoutingPool,
			Proxy:       t.Proxy}
		collector.Start()
	*/
	analyzer.StartAnalyzer()

	target := dataminer.Target{
		Stocks:      stocks}
	target.RegisterModuleDownloader(&downloader.JJCC{}).RegisterModuleDownloader(&downloader.GDTJ{})

	target.Start()
}
