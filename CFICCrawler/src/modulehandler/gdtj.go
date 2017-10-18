package modulehandler

import (
	"fmt"
	"os"
	"htmlparser"
	"errors"
)

// The file to handle the data of "gdtj" module

type GDTJ struct {
	Code string
	ID string
	DateList map[string]bool  //date:isDownloaded
	CurrentDate string
	Doc *htmlparser.HTMLDoc
}

const (
	GDTJ_LOCATION = "D:/Work/MyDemo/go/golang/CFICCrawler/resource/" //"E:/Programing/golang/CFICCrawler/resource/"
	GDTJ_HTML = "gdtj.html"
	GDTJ_QUARTER_LINK = "http://quote.cfi.cn/quote.aspx?stockid=%s&contenttype=gdtj&jzrq=%s"
)

func (gdtj *GDTJ)parseByDate(date string) error{
	file := GDTJ_LOCATION + gdtj.Code + ".html.modules/gdtj/" + GDTJ_HTML + "_" + date

	val,ok := gdtj.DateList[date]
	if !ok {
		return errors.New("Date does not exist")
	}

	if val { // The page has been downloaded, parse it directly.
		if doc, err := htmlparser.ParseFromFile(file); err != nil {
			return err
		} else {
			gdtj.Doc = doc
			gdtj.CurrentDate = date
		}
	} else {  // Download by url
		link := fmt.Sprintf(GDTJ_QUARTER_LINK, gdtj.ID, date)
		fmt.Println(link, file)

		if doc, err := gdtj.Doc.GDTJ_Request(link, file); err != nil {
			return err
		} else {
			gdtj.Doc = doc
			gdtj.CurrentDate = date
			gdtj.DateList[date] = true // Mark this page has been downloaded.
		}
	}

	return nil
}

func (gdtj *GDTJ)GetDateList() map[string]bool {
	if err := gdtj.getBasicData(); err != nil {
		return nil
	}

	return gdtj.DateList
}

// Get the shareholder in the specified perioid
func (gdtj *GDTJ)GetShareHolder(date string) ([]*htmlparser.ShareHolerInfo, error) {
	if err := gdtj.getBasicData(); err != nil {
		return []*htmlparser.ShareHolerInfo{}, err
	}

	if gdtj.CurrentDate != date {
		if err := gdtj.parseByDate(date); err != nil {
			return []*htmlparser.ShareHolerInfo{}, err
		}
	}

	// Get the free shareholders firstly. If they are not exist, get the major shareholders.
	sh := gdtj.Doc.GDTJ_GetShareholder(htmlparser.Free)
	if len(sh) == 0 {
		sh = gdtj.Doc.GDTJ_GetShareholder(htmlparser.Major)
	}

	return sh, nil
}

func (gdtj *GDTJ)getBasicData() error {
	if gdtj.ID == "" || len(gdtj.DateList) == 0  || gdtj.Doc == nil {
		file := GDTJ_LOCATION + gdtj.Code + ".html.modules/" + GDTJ_HTML
		doc, err := htmlparser.ParseFromFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse file %s faild, err:%s", file, err)
			return err
		}

		gdtj.ID = doc.GetStockId()
		gdtj.DateList = doc.GetDateList()
		gdtj.Doc = doc
		gdtj.CurrentDate = doc.GetCurrentDate()
	}

	return nil
}
