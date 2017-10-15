package modulehandler

import (
	"htmlparser"
	"fmt"
	"time"
	"strings"
)

// This file processes the history trade data from http://quotes.money.163.com/

// History Trade Data
type HTD struct {
	Code string
	Folder string // To write the history data

	Doc *htmlparser.HTMLDoc
}

/*
Example:
http://quotes.money.163.com/service/chddata.html?code=1000002&start=19910129&end=20161006&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP
code: 深市六位代码前加“1”，沪市股票代码前加“0”
start: 开始日期，如果想得到每只股票的所有历史交易数据，可以以公司上市日期来表达，8位数字，分别为yyyymmdd
end: 结束日期，表示的也是yyyymmdd八位数字
fields字段包括了开盘价、最高价、最低价、收盘价等
*/
const (
	// The Stock market started from 1990-12-19, so all the search start from this day.
	HTD_DOWNLOAD_LINK ="http://quotes.money.163.com/service/chddata.html?code=%s&start=19901219&end=%s&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCH"
)

func (htd *HTD)Download() error {
	var code string

	// 深市代码前加“1”，沪市股票代码前加“0”
	if strings.HasPrefix(htd.Code, "6") {
		code = "0" + htd.Code
	} else {
		code = "1" + htd.Code
	}

	year,month,day := time.Now().Date()
	link := fmt.Sprintf(HTD_DOWNLOAD_LINK, code, fmt.Sprintf("%d%d%d", year, month, day))
	file := htd.Folder + htd.Code + ".html.modules/htd/htd.csv"
	if err := htd.Doc.HTD_Request(link, file); err != nil {
		return err
	}

	return nil
}