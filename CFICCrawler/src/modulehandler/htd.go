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
	"sort"
)

// This file processes the history trade data from http://quotes.money.163.com/

// History Trade Data
type HTD struct {
	Code string
	Folder string // To write the stock history data
	SHMainIndexFile string // Where the Shang hai main index data file store in
	SZMainIndexFile string // Where the Shen zhen main index data file store in
	GEMfile	string	// Where the growth enterprises market data file store in

	Doc *htmlparser.HTMLDoc
}

type HTData struct {
	Date string			//历史数据日期
	Code string			//股票代码
	Name string			//股票名称
	ClosePrice float32  	//收盘价
	HighPrice float32   	// 最高价
	LowPrice float32    	//最低价
	StartPrice float32  	//开盘价
	PClosePrice float32 	//前收盘
	UDShortfall float32  //涨跌额
	UDRange	float32		//涨跌幅
	TurnoverRate float32	//换手率
	VOL float32			//成交量
	AMO float32			//成交金额
	TotalValue float32	//总市值
	FreeValue float32	// 流通市值
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

// logger
var logger = utility.GetLogger()

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
			data.ClosePrice = utility.String2Folat32(item)
		case 4:
			data.HighPrice = utility.String2Folat32(item)
		case 5:
			data.LowPrice = utility.String2Folat32(item)
		case 6:
			data.StartPrice = utility.String2Folat32(item)
		case 7:
			data.PClosePrice = utility.String2Folat32(item)
		case 8:
			data.UDShortfall = utility.String2Folat32(item)
		case 9:
			data.UDRange = utility.String2Folat32(item)
		case 10:
			data.TurnoverRate = utility.String2Folat32(item)
		case 11:
			data.VOL = utility.String2Folat32(item)
		case 12:
			data.AMO = utility.String2Folat32(item)
		case 13:
			data.TotalValue = utility.String2Folat32(item)
		case 14:
			data.FreeValue = utility.String2Folat32(item)
		}
	}

	return data
}

// Find the trade data by giving date list, and return them.
func (htd *HTD)getData(dateList []string) map[string]*HTData {
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

	var result = make(map[string]*HTData)
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
			//result = append(result, data)
			result[data.Date] = data
			cnt--

			logger.DEBUG(fmt.Sprintf("%s %s %s %f %f %f %f %f %f %f %f %f %f %f %f ",
				data.Name,
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
			))
		}
	}

	return result
}

// Get the Shang hai main index data by gaving date list
func (htd *HTD) getSHMainIndexdata(dateList []string) map[string]*HTData{
	// TODO
	var result = make(map[string]*HTData)

	return result
}

// Get the Shen zhen main index data by gaving date list
func (htd *HTD) getSZMainIndexdata(dateList []string) map[string]*HTData{
	// TODO
	var result = make(map[string]*HTData)

	return result
}

// Get the growth enterprises market data by gaving date list
func (htd *HTD) getGEMdata(dateList []string) map[string]*HTData{
	// TODO
	var result = make(map[string]*HTData)

	return result
}

// Check if the date is not in the weekend, if yes, change it to Friday.
func changeDate(dateList []string) []string{
	var result []string
	for _,item := range dateList {
		temp := item
		t, _ := time.Parse("2006-01-02", item)
		if t.Weekday().String() == "Sunday" {
			d, _ := time.ParseDuration("-48h")
			temp = t.Add(d).Format("2006-01-02")
			fmt.Println("changed Sunday", temp)
		}
		if t.Weekday().String() == "Saturday" {
			d, _ := time.ParseDuration("-24h")
			temp = t.Add(d).Format("2006-01-02")
			fmt.Println("changed Saturday", temp)
		}
		result = append(result, temp)
	}

	return result
}

func (htd *HTD)Analyse(focusSHIs map[string][]*htmlparser.ShareHolerInfo) {
	// To save all the date for every one focus shareholder
	shDateMap := make(map[string][]string)
	for fundname, shilist := range focusSHIs {
		for _, shi := range shilist {
			shDateMap[fundname] = append(shDateMap[fundname], shi.Date)
		}
	}

	for fundname, dateList := range shDateMap {
		//"D:/Work/MyDemo/go/golang/CFICCrawler/resource/"
		filename := "E:/Programing/golang/CFICCrawler/resource/" + fundname + "/" + htd.Code + ".csv"
		utility.WriteToFile(filename, "Date,Count,Ratio,Price,SH,SZ,GEM")

		dlist := changeDate(dateList)
		// Sort the date list so that searching data is sorted.
		sort.Strings(dlist)

		// Get the stock, main index and GEM history data, and write to file.
		// So that the history data can compare together.
		mapStockHistoryData := htd.getData(dlist)
		//mapSHMainIndexHistoryData := htd.getSHMainIndexdata(dateList)
		//mapSZMainIndexHistoryData := htd.getSZMainIndexdata(dateList)
		//mapGEMHistoryData := htd.getGEMdata(dateList)

		for _, date := range dlist {
			//if dayData, ok := mapStockHistoryData[date]; ok {
				for _, shi := range focusSHIs[fundname] {
					if shi.Date == date {

						line := fmt.Sprintf("%s,%d,%f,%f",
											shi.Date,
											shi.Count,
											shi.Ratio,
											mapStockHistoryData[date].StartPrice,
											)
						/*if data,ok := mapSHMainIndexHistoryData[date];ok {
							line = fmt.Sprintf("%s,%f", data.StartPrice)
						}
						if data,ok := mapSZMainIndexHistoryData[date];ok {
							line = fmt.Sprintf("%s,%f", data.StartPrice)
						}
						if data,ok := mapGEMHistoryData[date];ok {
							line = fmt.Sprintf("%s,%f", data.StartPrice)
						}*/

						logger.DEBUG(line)
						utility.WriteToFile(filename, line)
					}
				}
			//}
		}
	}
}