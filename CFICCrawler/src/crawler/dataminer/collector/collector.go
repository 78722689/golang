package collector

import (
	"htmlparser"
	"utility"
)

func StartCollectJJCC(file string) error {
	doc, err := htmlparser.ParseFromFile(file)
	if err != nil {
		utility.GetLogger().Errorf("Parse file failure, %s", err)
		return err
	}

	doc.JJCC_GetJJCCData()

	return nil
}