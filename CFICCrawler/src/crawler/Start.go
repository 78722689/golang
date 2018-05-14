package crawler

import (
	"crawler/dataminer"
	"routingpool"
	"crawler/dataminer/downloader"
)

func StartCrawl(pool *routingpool.ThreadPool, stocks []string) {

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
		RoutingPool: pool,
		Stocks:      stocks}

	target.RegisterModuleDownloader(&downloader.JJCC{}).RegisterModuleDownloader(&downloader.GDTJ{})

	target.Start()
}
