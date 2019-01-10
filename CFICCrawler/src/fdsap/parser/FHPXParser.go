package parser

import (
	"utility"
	"strings"
)

type FHPX_DATA struct {
	ExDividendDate string
	TransformNum int8		// 10 转 ?
	OfferNum int8			// 10 送 ?
	BTaxCashDividend float32	// 10 派 ？(税前)
	ATaxCashDividend float32	// 10 派 ？(税后)
}


func (doc *HTMLDoc)GetFHPXData() []*FHPX_DATA {
	var result []*FHPX_DATA

	doc.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if table.Nodes[0].GetAttrByName("id") != "tabh" {
			return
		}

		table.Find(TagNode, "tr").Each(func(i int, tr *Selection) {
			if i <= 1 {
				return
			}

			data := FHPX_DATA{}
			tr.Find(TagNode, "td").Each(func(i int, td *Selection) {
				td.Find(TextNode, "").Each(func(_ int, tn *Selection) {

					value := strings.TrimSpace(tn.Nodes[0].Root.Data)
					switch i {
					case 4:
						data.OfferNum = int8(utility.String2Int64(value))
					case 5:
						data.TransformNum = int8(utility.String2Int64(value))
					case 6:
						data.BTaxCashDividend = utility.String2Folat32(value)
					case 7:
						data.ATaxCashDividend = utility.String2Folat32(value)
					case 8:
						data.ExDividendDate = value
					}
				})
			})

			if data.ExDividendDate != "--" {
				result = append(result, &data)
			}
		})
	})

	return result
}