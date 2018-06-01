package collector

import (
	"htmlparser"
	"utility"
	"crawler/dataminer/analyzer"
)

func CollectJJCC(file string) error {
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		utility.GetLogger().Errorf("Parse file failure, %s", err)
		return err
	}

	analyzer.PutData(doc.JJCC_GetJJCCData())

	return nil
}