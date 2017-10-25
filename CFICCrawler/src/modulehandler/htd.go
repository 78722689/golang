package modulehandler

import (
	"htmlparser"
	"fmt"
	"time"
	"strings"
	"os"
	"bufio"

	"github.com/axgle/mahonia"
	"io"
	"utility"
)

// This file processes the history trade data from http://quotes.money.163.com/

// History Trade Data
type HTD struct {
	Code string
	Folder string // To write the history data

	Doc *htmlparser.HTMLDoc
}

type HTData struct {
	Date string			//历史数据日期
	Code string			//股票代码
	Name string			//股票名称
	ClosePrice string  	//收盘价
	HighPrice string   	// 最高价
	LowPrice string    	//最低价
	StartPrice string  	//开盘价
	PClosePrice string 	//前收盘
	UDShortfall string  //涨跌额
	UDRange	string		//涨跌幅
	TurnoverRate string	//换手率
	VOL string			//成交量
	AMO string			//成交金额
	TotalValue string	//总市值
	FreeValue string	// 流通市值
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
	HTD_DOWNLOAD_LINK ="http://quotes.money.163.com/service/chddata.html?code=%s&start=19901219&end=%s&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
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



func (htd *HTD)convert2HTData(line string) *HTData{
	data := &HTData{}

	for index, item := range strings.Split(strings.TrimSpace(line), ",") {
		switch index {
		case 0:
			data.Date = item
		case 1:
			data.Code = strings.Trim(item, "'")
		case 2:
			data.Name = item
		case 3:
			data.ClosePrice = item
		case 4:
			data.HighPrice = item
		case 5:
			data.LowPrice = item
		case 6:
			data.StartPrice = item
		case 7:
			data.PClosePrice = item
		case 8:
			data.UDShortfall = item
		case 9:
			data.UDRange = item
		case 10:
			data.TurnoverRate = item
		case 11:
			data.VOL = item
		case 12:
			data.AMO = item
		case 13:
			data.TotalValue = item
		case 14:
			data.FreeValue = item
		}
	}

	return data
}

// Find the trade data by giving date list, and return them.
func (htd *HTD)getData(dateList []interface{}) []*HTData{
	file := htd.Folder + htd.Code + ".html.modules/htd/htd.csv"

	f, err := os.Open(file)
	if  err != nil {
		fmt.Fprintf(os.Stderr, "Open file failure. %s\n", file)
		return nil
	}

	// Decode data due to Chinese
	decoder := mahonia.NewDecoder("gbk")
	reader := bufio.NewReader(decoder.NewReader(f))

	cnt := len(dateList)

	var result []*HTData
	// Loop the file line by line
	for {
		if cnt == 0 {
			break
		}

		l, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Fprintf(os.Stderr, "Read file failure. %s\n", file)
			break
		}

		data := htd.convert2HTData(l)
		if data.Date != "" && utility.Contains(dateList, data.Date) {
			result = append(result, data)
			cnt--

			fmt.Println(data.Name,
				data.Code,
				data.Date,
				data.UDRange,
				data.UDShortfall,
				data.PClosePrice,
				data.StartPrice,
				data.LowPrice,
				data.HighPrice,
				data.ClosePrice,
				data.AMO,
				data.FreeValue,
				data.TotalValue,
				data.TurnoverRate,
				data.VOL,
			)
		}
	}

	return result
}

func (htd *HTD)Analyse(dateList []interface{}) {
	htd.getData(dateList)

}