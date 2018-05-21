package crawler

import (
	"crawler/dataminer"
	"crawler/dataminer/downloader"
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
	//dataminer.StartAnalyse(codeChannel, pool)

	target := dataminer.Target{
		Stocks:      stocks}

	target.RegisterModuleDownloader(&downloader.JJCC{}).RegisterModuleDownloader(&downloader.GDTJ{})

	target.Start()
}
