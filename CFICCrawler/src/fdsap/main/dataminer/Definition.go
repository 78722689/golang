package dataminer

import (
	"fdsap/crawler/dataminer/downloader"
)

type Target struct {
	Stocks    []string // Empty for downloading all stocks.
	Modules   []downloader.Moduler
}

