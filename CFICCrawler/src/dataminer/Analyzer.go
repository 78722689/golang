package dataminer

import (
	"httpcontroller"
	"routingpool"
	"utility"
)

var logger = utility.GetLogger()

// The taget to analyze
type Target struct {
	Code        string
	Folder      string
	RoutingPool *routingpool.ThreadPool
	SyncChan    chan bool

	Proxy *httpcontroller.Proxy
}

func (t *Target) Start() {
	c := CollectInfo{Code: t.Code, Folder: t.Folder, Proxy: t.Proxy}
	c.StartMonitorDownloadStatus(t.RoutingPool, t.SyncChan)

}
