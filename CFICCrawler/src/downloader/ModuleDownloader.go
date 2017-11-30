package downloader

import (
	"sync"
	"httpcontroller"
	"htmlparser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"utility"
	"time"
	"routingpool"
)

var logger = utility.GetLogger()

const (
	STOCK_LIST_URL string = "http://quote.cfi.cn/stockList.aspx?t=11"
	QUOTE_HOMEPAGE_URL string = "http://quote.cfi.cn/"
)

type DownloadInfo struct {
	Foler string
	Proxy *httpcontroller.Proxy
	Overwrite bool

	RoutingPool *routingpool.ThreadPool
}

type DownloadTask struct {
	Name string
}

func (info *DownloadInfo)Task(id int) {
	//fmt.Println(fmt.Sprintf("Thread - %d is running with DownloadTask - %s", id, task.Name))
	time.Sleep(time.Second*2)
}

func (d *DownloadInfo)DownloadAll() {
	d.downloadStocksHomePage([]string{}, d.Overwrite)
	//d.downloadModules([]string{}, false)
}

func (d *DownloadInfo)DownloadByStockIDs(ids []string) {
	d.downloadStocksHomePage(ids, d.Overwrite)
	//d.downloadModules([]string{}, true)
}

// Download the stocks home page, if overwrite is true, the exist home page will be rewrite.
func (d *DownloadInfo)downloadStocksHomePage(ids []string, overwrite bool) {
	// Request homepage to get all the stocks
	request := httpcontroller.Request {
		Url   : STOCK_LIST_URL,
		Proxy : d.Proxy,
	}
	root,_ := request.Get()

	doc,err := htmlparser.ParseFromNode(root)
	if err != nil {
		logger.ERROR(fmt.Sprintf("Parse file error, %v", err))

		os.Exit(1)
	}

	mainTask := routingpool.Caller{Name:"Main-download-task", Call: func(id int) {
		fmt.Println("in caller...", ids)
		time.Sleep(time.Second*3)
		fmt.Println("in caller...after sleep", ids)

		for _, stockinfo := range doc.GetStocks(ids) {
			fmt.Println("stockinfo caller....")
			c := func (id int) {
				logger.INFO(fmt.Sprintf("Downloading link:%v name:%v, number:%v\r\n", stockinfo.Link, stockinfo.Name, stockinfo.Number))

				stock_request := httpcontroller.Request {
					Url  : QUOTE_HOMEPAGE_URL + stockinfo.Link,
					File : d.Foler + stockinfo.Number + "/" + stockinfo.Link,
					Proxy: d.Proxy,
					OverWrite:overwrite}
				stock_request.Get()
			}
			d.RoutingPool.PutTask(routingpool.Caller{Name:"Download-Homepage",Call:c})
		}
		fmt.Println("end caller...", ids)

	}}

	//task := routingpool.Caller{Name:"DownloaderHomepage",Call:caller}
	d.RoutingPool.PutTask(mainTask)
}

func (d *DownloadInfo)downloadModules(ids []string, overwrite bool) {
	wg := sync.WaitGroup{}

	// To walk the folder in order to find out the stock homepage html.
	err := filepath.Walk(d.Foler, func(path string, fi os.FileInfo, err error) error {
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

		wg.Add(1)
		// Begin to request the modules
		// filters: specify only to reqeust the interesting modules.
		go func(urls []string, filters []string) {

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

					file := d.Foler + values[len(values)-1] + ".modules/" + filename + ".html"
					request := httpcontroller.Request{
						Proxy: d.Proxy,
						Url:  QUOTE_HOMEPAGE_URL + url,
						File: file,
						OverWrite:overwrite}
					request.Get()
				}
			}

		}(stock_modules_url, filters)

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "No any file found in folder %s, err:%s", d.Foler, err)
	}

	wg.Wait()
}
