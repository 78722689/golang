package analyzer

import (
	"htmlparser"
	"httpcontroller"
	"modulehandler"
	"routingpool"
	"github.com/spf13/viper"
)

type stAnalysisRunner struct {
	code         string
	sourcefolder string
	wait         chan bool
	proxy        *httpcontroller.Proxy
}

func (r *stAnalysisRunner) caller(id int) {
	<-r.wait

	gdtj := &modulehandler.GDTJ{Code: r.code, Folder: r.sourcefolder, Proxy: r.proxy}
	targetFunds := []string{"香港中央结算有限公司"} //"上海汽车集团股份有限公司", "深圳市楚源投资发展有限公司"
	result := make(map[string][]*htmlparser.ShareHolerInfo)

	// Find out the fund if it is in the reporter
	for key, _ := range gdtj.GetDateList() {
		if shList, err := gdtj.GetShareHolder(key); err == nil {
			for _, fundname := range targetFunds {
				for _, sh := range shList {
					if sh.Name == fundname {
						result[fundname] = append(result[fundname], sh)

						//logger.Debugf("Found %s in %s", fundname, key)
						break
					}
				}
			}
		}
	}

	//htd.GetFundsPerformance(result, proxy)
}