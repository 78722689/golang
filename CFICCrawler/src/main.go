package main

import (
	"fmt"
	"modulehandler"
	"htmlparser"
	"utility"
	"httpcontroller"
)

// proxy //http://203.17.66.133:8000   http://203.17.66.134:8000
func main() {
	// Log setting
	logger := utility.GetLogger()
	logger.SetMinorLogLevel(utility.DEBUG)

	//var proxy *httpcontroller.Proxy = nil
	proxy := &httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"}
	folder := "D:/Work/MyDemo/go/golang/CFICCrawler/resource/"
	//folder := "E:/Programing/golang/CFICCrawler/resource/"

	/*
	downloader := DownloadInfo{Foler:folder, Proxy:proxy}
	downloader.Download()
	*/

	code := "601700" //"601699"
	gdtj := &modulehandler.GDTJ{Code:code, Folder:folder}
/*
	if sh, err := gdtj.GetShareHolder("2015-12-31", proxy); err == nil {
		for _, item:= range sh {
			fmt.Fprintf(os.Stdout, "Name:%s, Count:%s, Ratio:%s\n", item.Name, item.Count, item.Ratio)
		}
	}
*/
	htd := modulehandler.HTD{Code : code,
							Folder : folder}

	htd.Download(proxy)

	funds := []string{"全国社保基金一零四组合","中国工商银行-嘉实策略增长混合型证券投资基金"}
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

	// Requirements:
	// 1. 加入同时期大盘指数走势----->done
	// 2. 计算除权价格？
	// 3. 计算分红数据，持股变动后盈利以及总盈利。
	htd.GetFundsPerformance(result, proxy)

	logger.DEBUG("main is end...........................")

}
