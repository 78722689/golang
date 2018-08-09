package collector

import (
	"htmlparser"
	"utility"
	"crawler/dataminer/analyzer"
)

func CollectJJCC(stocknumber, stockname, file, recordDate string) error {
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		utility.GetLogger().Errorf("Parse file failure, %s", err)
		return err
	}
	analyzer.PushDataIntoRedis(doc.JJCC_GetJJCCData(stocknumber, stockname, recordDate))

	return nil
}
