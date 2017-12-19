package main

import (
	"utility"
	//	"sync"
	"downloader"
	//"fmt"
	//"htmlparser"
	"httpcontroller"
	//"modulehandler"
	"routingpool"
	"runtime"
)

// proxy //http://203.17.66.133:8000   http://203.17.66.134:8000
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Log setting
	logger := utility.GetLogger()
	logger.SetMinorLogLevel(utility.DEBUG)

	pool := routingpool.GetPool(100, 100)
	pool.Start()
	/*
		for i := 0; i<=20; i++ {
			go func(id int) {
				for j := 0; j<=20; j++ {
					download := downloader.DownloadTask{Name:fmt.Sprintf(" Task{id - %d, queue-%d}", id, j)}
					task := &routingpool.MyTask{Name:"customer-task", Call:download.Task}
					myPool.PutTask(task)
				}
			}(i)
		}
	*/

	var proxy *httpcontroller.Proxy = nil
	//proxy := &httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"}
	//folder := "D:/Work/MyDemo/go/golang/CFICCrawler/resource/download/"
	folder := "E:/Programing/golang/CFICCrawler/resource/download/"

	downloader := downloader.DownloadInfo{Folder: folder, Proxy: proxy, Overwrite: false, RoutingPool: pool}
	downloader.DownloadByStockIDs([]string{"600089"}) //, "600096", "600036"

	// Waiting for all threads finish and exit
	pool.Wait()

	//code := "600036" //"600089"
	//gdtj := &modulehandler.GDTJ{Code: code, Folder: folder, Proxy : proxy}

	/*
		if sh, err := gdtj.GetShareHolder("2015-12-31"); err == nil {
			for _, item:= range sh {
				fmt.Fprintf(os.Stdout, "Name:%s, Count:%s, Ratio:%s\n", item.Name, item.Count, item.Ratio)
			}
		}
	*/

	/*htd := modulehandler.HTD{Code: code, Folder: folder, Proxy: proxy}
	htd.Download()

	funds := []string{"上海汽车集团股份有限公司", "香港中央结算有限公司", "深圳市楚源投资发展有限公司"}
	result := make(map[string][]*htmlparser.ShareHolerInfo)

	// Find out the fund if it is in the reporter
	for key, _ := range gdtj.GetDateList() {
		if shList, err := gdtj.GetShareHolder(key, proxy); err == nil {
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

	htd.GetFundsPerformance(result, proxy)
	*/
	logger.DEBUG("main is end...........................")

}
