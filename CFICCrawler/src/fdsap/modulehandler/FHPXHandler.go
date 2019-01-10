package modulehandler

import (
	"htmlparser"
)

type FHPX_INFO struct {
	Code   string
	Folder string

	Doc *htmlparser.HTMLDoc
}

const (
	FHPX_HOMEPAGE = "fhpx.html"
)

func (fhpx *FHPX_INFO) GetFHPXData() ([]*htmlparser.FHPX_DATA, error) {
	path := fhpx.Folder + fhpx.Code + "/modules/" + FHPX_HOMEPAGE

	doc, err := htmlparser.ParseFromFile(path)
	if err != nil {
		logger.Errorf("Parse file faile, %s", err)
		return nil, err
	}

	data := doc.GetFHPXData()
	for _, d := range data {
		logger.Debugf("ExDividendDate:%s, BTaxCashDividend:%f, TransformNum:%d, ATaxCashDividend:%f OfferNum:%d",
			d.ExDividendDate,
			d.BTaxCashDividend,
			d.TransformNum,
			d.ATaxCashDividend,
			d.OfferNum)
	}

	return data, nil
}
