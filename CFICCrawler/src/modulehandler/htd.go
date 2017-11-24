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
	"httpcontroller"
)

// This file processes the history trade data from http://quotes.money.163.com/

// History Trade Data
type HTD struct {
	Code string
	Folder string // To write the stock history data
	StockDataFile	string // Where the stock history data file store in
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

type Main_Index_Type uint8
const  (
	Stock Main_Index_Type = iota  // Stock
	SH 	// Shanghai main index
	SZ	// Shenzhen main index
	GEM	// growth enterprises market
)

// logger
var logger = utility.GetLogger()

func (htd *HTD)Download(proxy *httpcontroller.Proxy) error {
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
	if err := htd.Doc.HTD_Request(link, file, proxy); err != nil {
		return err
	}
	htd.StockDataFile = file

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
func (htd *HTD)getData(dateList []interface{}, mit Main_Index_Type) map[string]*HTData {
	file := ""
	switch mit {
	case Stock:
		file = htd.Folder + htd.Code + ".html.modules/htd/htd.csv"
	case SH:
		file = htd.SHMainIndexFile
	case SZ:
		file = htd.SZMainIndexFile
	case GEM:
		file = htd.GEMfile
	}

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
func (htd *HTD) getSHMainIndexdata(dateList []interface{}, proxy *httpcontroller.Proxy) map[string]*HTData{
	if htd.SHMainIndexFile == "" {
		year,month,day := time.Now().Date()
		link := fmt.Sprintf(HTD_DOWNLOAD_LINK, "0000001", fmt.Sprintf("%d%d%d", year, month, day))
		file := htd.Folder + "000001/htd/htd.csv"
		if err := htd.Doc.HTD_Request(link, file, proxy); err != nil {
			logger.ERROR(fmt.Sprintf("Fetch Shang hai main index data failure, %s", link))
			return nil
		}
		htd.SHMainIndexFile = file
	}
	return htd.getData(dateList, SH)
}

// Get the Shen zhen main index data by gaving date list
func (htd *HTD) getSZMainIndexdata(dateList []interface{}, proxy *httpcontroller.Proxy) map[string]*HTData{
	if htd.SZMainIndexFile == "" {
		year,month,day := time.Now().Date()
		link := fmt.Sprintf(HTD_DOWNLOAD_LINK, "1399001", fmt.Sprintf("%d%d%d", year, month, day))
		file := htd.Folder + "399001/htd/htd.csv"
		if err := htd.Doc.HTD_Request(link, file, proxy); err != nil {
			logger.ERROR(fmt.Sprintf("Fetch Shen zhen main index data failure, %s", link))
			return nil
		}
		htd.SZMainIndexFile = file
	}
	return htd.getData(dateList, SZ)
}

// Get the growth enterprises market data by gaving date list
func (htd *HTD) getGEMdata(dateList []interface{}, proxy *httpcontroller.Proxy) map[string]*HTData{
	if htd.GEMfile == "" {
		year,month,day := time.Now().Date()
		link := fmt.Sprintf(HTD_DOWNLOAD_LINK, "1399006", fmt.Sprintf("%d%d%d", year, month, day))
		file := htd.Folder + "399006/htd/htd.csv"
		if err := htd.Doc.HTD_Request(link, file, proxy); err != nil {
			logger.ERROR(fmt.Sprintf("Fetch Shen zhen main index data failure, %s", link))
			return nil
		}
		htd.GEMfile = file
	}
	return htd.getData(dateList, GEM)
}

// Check if the date is not in the weekend, if yes, change it to Friday.
// Return origin date list and changed date list.
func getNoWeekendDateList(dateList []string) (map[string]string) {
	var result = make(map[string]string)
	for _,item := range dateList {
		temp := item
		t, _ := time.Parse("2006-01-02", item)
		if t.Weekday().String() == "Sunday" {
			d, _ := time.ParseDuration("-48h")
			temp = t.Add(d).Format("2006-01-02")
			logger.WARN(fmt.Sprintf("Changed date Sunday (%s) to Friday (%s)", item, temp))
		}
		if t.Weekday().String() == "Saturday" {
			d, _ := time.ParseDuration("-24h")
			temp = t.Add(d).Format("2006-01-02")
			logger.WARN(fmt.Sprintf("Changed date Saturday (%s) to to Friday (%s)", item, temp))
		}
		result[item] = temp
	}

	return result
}

// Get the nearest FHPX data by date
func (htd *HTD)getNearestFHPXDataByDate(date string) *htmlparser.FHPX_DATA {
	// Get all the FHPX data on the stock
	fhpxInfo := FHPX_INFO{Code : htd.Code, Folder : htd.Folder}
	fhpxDatalist,err := fhpxInfo.GetFHPXData()
	if err != nil {
		logger.ERROR("Get FHPX data failure")
		return nil
	}

	htdDate := utility.String2Date(date)
	var result *htmlparser.FHPX_DATA
	for _, fhpx := range fhpxDatalist {
		exDividendDate := utility.String2Date(fhpx.ExDividendDate)
		if      (htdDate.Year() - exDividendDate.Year()) == 0 &&
			((htdDate.Month() - exDividendDate.Month()) >=0 &&
				(htdDate.Month() - exDividendDate.Month()) <= 3) {
			result = fhpx

			logger.DEBUG(fmt.Sprintf("FouND FHPX data, date-%s vs htd date %s",
				fhpx.ExDividendDate,
				date))
		}
	}

	return result
}

//func (htd *HTD)getMainIndex

func (htd *HTD)Analyse(focusSHIs map[string][]*htmlparser.ShareHolerInfo, proxy *httpcontroller.Proxy) {
	// To save all the date for every one focus shareholder
	shDateMap := make(map[string][]string)
	for fundname, shilist := range focusSHIs {
		for _, shi := range shilist {
			shDateMap[fundname] = append(shDateMap[fundname], shi.Date)
		}
	}

	for fundname, dateList := range shDateMap {
		filename := htd.Folder + fundname + "/" + htd.Code + ".csv"
		utility.WriteToFile(filename, "Date,Count,Ratio,Price,SH,SZ,GEM,ExDividendPrice,FHPXProfit,ChangeProfit,OfferNumber,TransformNumber")

		// Sort the date list so that searching data is sorted.
		sort.Strings(dateList)

		// Get the date without weekend.
		noWeekendDatemap := getNoWeekendDateList(dateList)


		dlist := utility.Values(noWeekendDatemap)
		// Get the stock, main index and GEM history data, and write to file.
		// So that the history data can compare together.
		mapStockHistoryData := htd.getData(dlist, Stock)
		mapSHMainIndexHistoryData := htd.getSHMainIndexdata(dlist, proxy)
		mapSZMainIndexHistoryData := htd.getSZMainIndexdata(dlist, proxy)
		mapGEMHistoryData := htd.getGEMdata(dlist, proxy)

		preCount := 0.0
		for _, date := range dateList {
			for _, shi := range focusSHIs[fundname] {
				if shi.Date == date {
					if data,ok := mapStockHistoryData[noWeekendDatemap[date]]; ok {
						FHPXProfit := 0		// The profit in ex-dividend day
						ChangeProfit := 0   // The profit when count happened change.
						changeCount := preCount - utility.String2Folat64(shi.Count)
						if changeCount != 0.0 {
							ChangeProfit = int(changeCount * 10000 * float64(data.StartPrice))
							logger.DEBUG(fmt.Sprintf("date-%s preCount-%f, thisCount-%s, ChangeProfit-%f", shi.Date, preCount, shi.Count, ChangeProfit))
						}

						fhpxData := htd.getNearestFHPXDataByDate(date)
						// Caculate the profit on the reporting day
						priceOnFHPXDay := 0.0
						offerNum := 0
						transformNum := 0
						if fhpxData != nil {
							mapTemp := htd.getData([]interface{}{fhpxData.ExDividendDate}, Stock)
							priceOnFHPXDay = float64(mapTemp[fhpxData.ExDividendDate].StartPrice)
							FHPXProfit = int(utility.String2Folat64(shi.Count) * (10000/10) * float64(fhpxData.ATaxCashDividend))
							offerNum = int(fhpxData.OfferNum)
							transformNum = int(fhpxData.TransformNum)
						}

						SHStartPrice := 0.0
						SZStartPrice := 0.0
						GEMStartPrice := 0.0
						if mapSHMainIndexHistoryData[noWeekendDatemap[date]] != nil {
							SHStartPrice= float64(mapSHMainIndexHistoryData[noWeekendDatemap[date]].StartPrice)
						}
						if mapSZMainIndexHistoryData[noWeekendDatemap[date]] != nil {
							SZStartPrice= float64(mapSZMainIndexHistoryData[noWeekendDatemap[date]].StartPrice)
						}
						if mapGEMHistoryData[noWeekendDatemap[date]] != nil {
							GEMStartPrice= float64(mapGEMHistoryData[noWeekendDatemap[date]].StartPrice)
						}

						line := fmt.Sprintf("%s,%s,%f,%f,%f,%f,%f,%f,%d,%d,%d,%d",
							shi.Date,
							shi.Count,
							shi.Ratio,
							data.StartPrice,
							SHStartPrice,
							SZStartPrice,
							GEMStartPrice,
							priceOnFHPXDay,
							FHPXProfit,
							ChangeProfit,
							offerNum,
							transformNum,
						)

						logger.DEBUG(line)
						utility.WriteToFile(filename, line)

						preCount = utility.String2Folat64(shi.Count)
					}
				}
			}
		}
	}
}