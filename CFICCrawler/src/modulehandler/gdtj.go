package modulehandler

import (
	"fmt"
	"os"
	"htmlparser"

)

// The file to handle the data of "gdtj" module

const (
	GDTJ_LOCATION = "E:/Programing/GO/CFICCrawler/resource/" //"D:/Work/MyDemo/go/golang/CFICCrawler/resource/"
	GDTJ_HTML = "gdtj.html"
	GDTJ_QUARTER_LINK = "http://quote.cfi.cn/quote.aspx?stockid=%s&contenttype=gdtj&jzrq=%s"
)

type GDTJ_INFO struct {
	StockID string
	fileLocation string
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
	}

	for _, shi := range doc.GetShareholder(htmlparser.Free) {
		fmt.Fprintf(os.Stdout, "name: %s count:%s, ratio:%s\n", shi.Name, shi.Count, shi.Ratio)
	}

	return nil
}
