package main

import (
	"github.com/op/go-logging"
	"fdsap/http"
	"fdsap/utility"
	"fdsap/main/dataminer"
	"fdsap/main/dataminer/downloader"
)

type Context struct {
	Log  *logging.Logger
	HttpHandler *http.HttpHandler
}

var context *Context

func init() {
	context = &Context{}
}

func (c *Context) WaitAll() {
	c.HttpHandler.Wait()
}


func StartCrawl(stocks []string) {

	logger := utility.GetLogger()
	context.Log = logger
	context.HttpHandler,_ = http.NewHttpHandler(&http.HttpHandlerConfig{128, 128, logger})

	context.WaitAll()
	/*
		collector := Collect{
			Code:        tempStockinfo.Number,
			Folder:      t.Folder,
			SyncChan:    syncChan,
			RoutingPool: t.RoutingPool,
			Proxy:       t.Proxy}
		collector.Start()
*/
	if false {
		//database.StartAnalyzer()

		target := dataminer.Target{
			Stocks: stocks}
		target.RegisterModuleDownloader(&downloader.JJCC{}).RegisterModuleDownloader(&downloader.GDTJ{}).RegisterModuleDownloader(&downloader.DOMAIN{})

		target.Start()
	}

	//analyser.StartAnalyse()
}
