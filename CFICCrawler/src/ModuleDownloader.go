package main

import (
	"sync"
	"httpcontroller"
	"htmlparser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

const (
	STOCK_LIST_URL string = "http://quote.cfi.cn/stockList.aspx?t=11"
	QUOTE_HOMEPAGE_URL string = "http://quote.cfi.cn/"
)

type DownloadInfo struct {
	Foler string
	Proxy *httpcontroller.Proxy
}

func (d *DownloadInfo)InitGoRouting() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (d *DownloadInfo)Download() {
	d.InitGoRouting()
	d.downloadHomePage()
	d.downloadModules()
}

func (d *DownloadInfo)downloadHomePage() {
	// Request homepage to get all the stocks
	request := httpcontroller.Request {
		Url   : STOCK_LIST_URL,
		Proxy : d.Proxy,
	}
	root,_ := request.Get()

	doc,err := htmlparser.ParseFromNode(root)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Parse file error, %v", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	// Download all the homepage of stocks and write to file.
	for _, stockinfo := range doc.GetAllStocks() {
		go func(si htmlparser.StockInfo) {
			fmt.Fprintf(os.Stdout, "Downloading link:%v name:%v, number:%v\r\n", si.Link, si.Name, si.Number)
			wg.Add(1)
			stock_request := httpcontroller.Request {
				Url  : QUOTE_HOMEPAGE_URL + si.Link,
				File : d.Foler + si.Link,
				Proxy: d.Proxy,
			}
			stock_request.Get()
			wg.Done()
		}(stockinfo)
	}

	wg.Wait()
}

func (d *DownloadInfo)downloadModules() {
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

					file := d.Foler + values[len(values)-1] + ".modules/" + filename + ".html"
					request := httpcontroller.Request{
						Proxy: d.Proxy,
						Url:  QUOTE_HOMEPAGE_URL + url,
						File: file,
					}

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
