package dataminer

import (
	"routingpool"
	"crawler/dataminer/downloader"
)

type Target struct {
	Stocks    []string // Empty for downloading all stocks.
	Modules   []downloader.Moduler
	RoutingPool *routingpool.ThreadPool
}

