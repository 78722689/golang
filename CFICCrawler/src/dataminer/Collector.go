package dataminer

import (
	"fmt"
	"htmlparser"
	"httpcontroller"
	"modulehandler"
	"routingpool"
)

type CollectInfo struct {
	Code   string
	Folder string

	Proxy *httpcontroller.Proxy
}

func (c *CollectInfo) StartMonitorDownloadStatus(routingPool *routingpool.ThreadPool, status chan bool) {
	runner := func(id int) {
		<-status

		gdtj := &modulehandler.GDTJ{Code: c.Code, Folder: c.Folder, Proxy: c.Proxy}
		/*if sh, err := gdtj.GetShareHolder("2015-12-31"); err == nil {
			for _, item := range sh {
				logger.DEBUG(fmt.Sprintf("Name:%s, Count:%s, Ratio:%f\n", item.Name, item.Count, item.Ratio))
			}
		}*/

		htd := modulehandler.HTD{Code: c.Code, Folder: c.Folder, Proxy: c.Proxy}
		htd.Download()

		funds := []string{"香港中央结算有限公司"} //"上海汽车集团股份有限公司", "深圳市楚源投资发展有限公司"
		result := make(map[string][]*htmlparser.ShareHolerInfo)

		// Find out the fund if it is in the reporter
		for key, _ := range gdtj.GetDateList() {
			if shList, err := gdtj.GetShareHolder(key); err == nil {
				for _, fundname := range funds {
					for _, sh := range shList {
						if sh.Name == fundname {
							result[fundname] = append(result[fundname], sh)

							logger.DEBUG(fmt.Sprintf("Found %s in %s", fundname, key))
							break
						}
					}
				}
			}
		}

		//htd.GetFundsPerformance(result, proxy)
	}

	routingPool.PutTask(routingpool.NewCaller("Collector", runner))
}
