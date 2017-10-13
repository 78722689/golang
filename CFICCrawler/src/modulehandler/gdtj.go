package modulehandler

import (
	"fmt"
	"os"
	"htmlparser"

	"httpcontroller"
)

// The file to handle the data of "gdtj" module

const (
	GDTJ_LOCATION = "D:/Work/MyDemo/go/golang/CFICCrawler/resource/" //"E:/Programing/GO/CFICCrawler/resource/"
	GDTJ_HTML = "gdtj.html"
	GDTJ_QUARTER_LINK = "http://quote.cfi.cn/quote.aspx?stockid=%s&contenttype=gdtj&jzrq=%s"
)

func ParseByDate(code string, id string, date string) {
	link := fmt.Sprintf(GDTJ_QUARTER_LINK, id, date)
	request := httpcontroller.Request{
		Proxy:&httpcontroller.Proxy{"HTTP", "10.144.1.10", "8080"},
		Url : link,
		File : GDTJ_LOCATION + code + ".html.modules/" + GDTJ_HTML + "_" + date,
	}

	request.Get()
}

func Parse(code string) error {
	file := GDTJ_LOCATION + code + ".html.modules/" + GDTJ_HTML
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse file %s faild, err:%s", file, err)
		return err
	}



	for _, d := range doc.GetDateList() {
		fmt.Println(d)
		ParseByDate(code, doc.GetStockId(), d)
	}
/*
	for _, shi := range doc.GetShareholder(htmlparser.Free) {
		fmt.Fprintf(os.Stdout, "name: %s count:%s, ratio:%s\n", shi.Name, shi.Count, shi.Ratio)
	}
*/
	return nil
}
