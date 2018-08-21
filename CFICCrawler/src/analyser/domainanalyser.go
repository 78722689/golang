package analyser

import (
	"github.com/garyburd/redigo/redis"
	"crawler/dataminer/database"
	"utility"
	"routingpool"
	"github.com/axgle/mahonia"
	"encoding/json"
	"time"
	"fmt"
	"strings"
)

var (
	logger = utility.GetLogger()
	redisPool *redis.Pool
	//domain_jjcc_chan chan map[string]map[string]map[string]float64
)

func init() {
	redisPool = database.RedisPool
	//domain_jjcc_chan = make(chan map[string]map[string]map[string]float64, 200)
}

// Implement Task interface
type caculator struct {
	routingpool.Base
}

func newCaculator(name string, call func(int)) *caculator {
	return &caculator{Base : routingpool.Base{Name: name, Call: call, Response: make(chan bool)}}
}

func (c *caculator) Run(id int) {
	c.Call(id)
}

/*
func all_total_jjcc_cacul() {
	for {
		select {
		case jjcc := <-domain_jjcc_chan:

		}
	}
}
*/
func DoDomainAnalyse() {
	if redisPool == nil {
		logger.Debug(" nil for redis connection")
	}

	domains, err := redis.Strings(redisPool.Get().Do("KEYS", "SET_DOMAIN_STOCKS_MAPPING_*"))
	if len(domains) == 0 || err != nil {
		logger.Errorf("Query Redis failed for DOMAINS, %s", err)
		return
	}
	all_domains_total_jjcc := make(map[string]map[string]map[string]float64)
	redisLangDecoder := mahonia.NewDecoder("gbk")
	for _, domain := range domains {
		decodedDomain := redisLangDecoder.ConvertString(domain)
		logger.Debugf("*********************Domain:%s*********************", decodedDomain)

		//domainCaculator := func(id int) {
		func(id int) {
			logger.Debugf("Thread-%d,Begin query doman&stocks mapping for domain %s*********************", id, decodedDomain)
			redisLangEncoder := mahonia.NewEncoder("gbk")
			stocks, err := redis.Strings(redisPool.Get().Do("SMEMBERS", redisLangEncoder.ConvertString(decodedDomain)))
			if err != nil {
				logger.Errorf("Thread-%d, Query Redis failed for doman&stocks mapping, domain: %s", id, decodedDomain)
			}
			totalJJCC := make(map[string]map[string]float64)  // {date:{count:0.1, value:0.1},}
			for _, stock := range stocks {
				logger.Debugf("Thread-%d, bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb stock:%s", id, stock)
				jjccRows, err := redis.Values(redisPool.Get().Do("LRANGE", "JJCC_" + stock, 0, -1))
				if err != nil {
					logger.Errorf("Thread-%d, Query Redis failed for JJCC for stock:%s", id, stock)
					continue
				}

				for _, row := range jjccRows {
					var jjccMap map[string]map[string]map[string]string
					json.Unmarshal(row.([]byte), &jjccMap)

					for date, jjdataMap := range jjccMap {
						logger.Debugf("Thread-%d, ================ date: %s =======================", id, date)
						totalCount := 0.0
						totalValue := 0.0
						for jjName, dataMap := range jjdataMap {
							logger.Debugf("Thread-%d, row JJName:%s, count:%s, value:%s", id, jjName, dataMap["count"], dataMap["value"])
							totalCount = totalCount + utility.String2Folat64(dataMap["count"])
							totalValue = totalValue + utility.String2Folat64(dataMap["value"])
						}
						logger.Debugf("Thread-%d, Total: count-%.4f, value-%.4f", id, totalCount, totalValue)

						// All JJCC data in one month plus to one record, because the stocks usually did not report JJCC in a day.
						t,_ := time.Parse("2006-01-02", date)
						totalDate := t.Format("2006-01")+"-01"
						if t.Month() == time.January || t.Month() == time.February || t.Month() == time.March {
							totalDate = t.Format("2006-")+"03-31"
						} else if t.Month() == time.April || t.Month() == time.May || t.Month() == time.June {
							totalDate = t.Format("2006-")+"06-30"
						} else if t.Month() == time.July || t.Month() == time.August || t.Month() == time.September {
							totalDate = t.Format("2006-")+"09-30"
						} else { //October November December
							totalDate = t.Format("2006-")+"12-31"
						}

						if t, ok := totalJJCC[totalDate]; ok {
							t["count"] = t["count"] + totalCount
							t["value"] = t["value"] + totalValue
						} else {
							totalJJCC[totalDate] = map[string]float64{"count":totalCount, "value":totalValue}
						}

						logger.Debugf("Thread-%d, =======================================================", id)
					}
				}

				logger.Debugf("Thread-%d, eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee %v", id, stock)
			}

			// Write Header
			//file := fmt.Sprintf("d:/out/result/%s.csv", strings.Replace(decodedDomain,"SET_DOMAIN_STOCKS_MAPPING_", "", 1))
			//utility.WriteToFile(file, "RecordDate,Count,Value")
			for date, jjcc := range totalJJCC {
				// Try to write result to file for eacch domain
				//line := fmt.Sprintf("%s,%.4f,%.4f", date, jjcc["count"], jjcc["value"])
				//utility.WriteToFile(file, line)

				// Caculating all domains JJCC
				if domain_total_jjcc, ok := all_domains_total_jjcc[date]; ok {
					domain_total_jjcc[decodedDomain] = jjcc
				} else {
					all_domains_total_jjcc[date] = map[string]map[string]float64{decodedDomain:jjcc}
				}
			}

			//domain_jjcc_chan <- map[string]map[string]map[string]float64{decodedDomain:totalJJCC}
		}(222)

		//routingpool.PutTask(newCaculator("Domain-Caculator", domainCaculator))
	}

	// Try to write all domains & all JJCC data to one file
	all_count_file := "d:/out/result/all_count.csv"
	all_value_file := "d:/out/result/all_value.csv"
	all_header := "RecordDate"

	for _, domain := range domains {
		decodedDomain := redisLangDecoder.ConvertString(strings.Replace(domain,"SET_DOMAIN_STOCKS_MAPPING_", "", 1))
		all_header = fmt.Sprintf("%s,%s", all_header, decodedDomain)
		fmt.Println(decodedDomain)
	}
	utility.WriteToFile(all_count_file, all_header)
	utility.WriteToFile(all_value_file, all_header)

	for date, total_jjcc :=range all_domains_total_jjcc {
		lineCount := date
		lineValue := date
		for _, domain := range domains {
			decodedDomain := redisLangDecoder.ConvertString(domain)
			//fmt.Println(decodedDomain)
			if jjcc, ok := total_jjcc[decodedDomain]; ok {
				lineCount = fmt.Sprintf("%s,%.4f", lineCount, jjcc["count"])
				lineValue = fmt.Sprintf("%s,%.4f", lineValue, jjcc["value"])
			} else {
				lineCount = fmt.Sprintf("%s,0", lineCount)
				lineValue = fmt.Sprintf("%s,0", lineValue)
			}

			//fmt.Println(lineCount)
			//fmt.Println(lineValue)
		}
		utility.WriteToFile(all_count_file, lineCount)
		utility.WriteToFile(all_value_file, lineValue)
	}
}