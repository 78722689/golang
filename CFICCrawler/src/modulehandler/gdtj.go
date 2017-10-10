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
)




func Parse(code string) error {
	file := GDTJ_LOCATION + code + ".html.modules/" + GDTJ_HTML
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse file %s faild, err:%s", file, err)
		return err
	}

	/*
	for _, d := range doc.GetDateList() {
		fmt.Println(d)
	}
*/
	doc.GetMajorShareholder()
	return nil
}
