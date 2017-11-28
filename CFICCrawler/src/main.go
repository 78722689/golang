package main

import (
	//"fmt"
	//"os"
	//"htmlparser"
	//"httpcontroller"
	"runtime"
	//"path/filepath"
	//"sync"
	//"httpcontroller"
	//"strings"

	//"modulehandler"
	"fmt"
	"modulehandler"
	"htmlparser"
	"utility"
	"httpcontroller"
)

const (
	//QUOTE_HOMEPAGE string = "http://quote.cfi.cn/"
	FOLDER_TOWRITE string = "E:/programing/GO/CFICCrawler/resource/"
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

	// goroutines settings
	runtime.GOMAXPROCS(runtime.NumCPU())
	/*
			wg := sync.WaitGroup{}

			// Request homepage to get all the stocks
			request := httpcontroller.Request {
				Url   : "http://quote.cfi.cn/stockList.aspx?t=11",
			}
			root,_ := request.Get()

			doc,err := htmlparser.ParseFromNode(root)
			if err != nil {
				fmt.Fprintf(os.Stdout, "Parse file error, %v", err)
				os.Exit(1)
			}

			// Download all the homepage of stocks and write to file.
			for _, stockinfo := range doc.GetAllStocks() {
				go func(si htmlparser.StockInfo) {
					fmt.Fprintf(os.Stdout, "Downloading link:%v name:%v, number:%v\r\n", si.Link, si.Name, si.Number)
					wg.Add(1)
					stock_request := httpcontroller.Request {
						Url  : QUOTE_HOMEPAGE + si.Link,
						File : FOLDER_TOWRITE + si.Link,
					}
					stock_request.Get()
					wg.Done()
				}(stockinfo)
			}

			wg.Wait()

		// To walk the folder in order to find out the stock homepage html.
		err := filepath.Walk(folder, func(path string, fi os.FileInfo, err error) error {
			strRet, _ := os.Getwd()
			ostype := os.Getenv("GOOS") // windows, linux

			if ostype == "windows" {
				strRet += "\\"
			} else if ostype == "linux" {
				strRet += "/"
			}

			if fi == nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}

			// Begin to parse the stock home page.
			fmt.Fprintf(os.Stdout, "Parsing file %s\n", path)
			doc, err := htmlparser.ParseFromFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Parse file %s faild, err:%s", path, err)
				return err
			}
			stock_modules_url := doc.GetModuleURL(fi.Name())

			// an example filter on a stock modules
			filters := []string{"gdtj", "fhpx"}

			// Begin to request the modules
			// filters: specify only to reqeust the interesting modules.
			go func(urls []string, filters []string) {
				wg.Add(1)
				defer wg.Done()
				for _, url := range urls {
					fmt.Fprintf(os.Stdout, "Checking stock module url:%s\n", url)

					var found bool = false
					for _, filter := range filters {
						if strings.Contains(url, filter) {
							found = true
							break
						}
					}

					if found {
						fmt.Fprintf(os.Stdout, "Found! goto request:%s\n", url)

						values := strings.Split(url, "/")
						var filename string
						if len(values) == 4 {
							filename = values[1] // example data "/tzzk/19770/601015.html"
						} else {
							filename = values[2] // example data "http://gg.cfi.cn/cbgg/19770/601015.html"
						}

						file := folder + values[len(values)-1] + ".modules/" + filename + ".html"
						request := httpcontroller.Request{
							Proxy: proxy,
							Url:  QUOTE_HOMEPAGE + url,
							File: file,
						}

						request.Get()
					}
				}

			}(stock_modules_url, filters)

			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "No any file found in folder %s, err:%s", FOLDER_TOWRITE, err)
		}

		wg.Wait()


*/
	code := "601699" //"601700"
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
	htd.Analyse(result, proxy)

	logger.DEBUG("main is end...........................")

}
