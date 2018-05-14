package htmlparser

import (
	"httpcontroller"
	"strings"
	"fmt"
	"utility"
	"time"
)

func (tree *HTMLDoc) JJCC_Request(url string, file string) (*HTMLDoc, error) {
	request := httpcontroller.Request{
		Url:       url,
		File:      file,
		OverWrite: false,
	}

	if _, err := request.Get(); err != nil {
		return nil, err
	}

	if doc, err := ParseFromFile(file); err != nil {
		return nil, err
	} else {
		return doc, nil
	}
}

func (tree *HTMLDoc) JJCC_GetJJCCData() []*ShareHolerInfo {
	var shiList []*ShareHolerInfo

	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {
			return
		}

		table.Find(TagNode, "tr").Each(func(i int, tr *Selection) {
			if i <= 1 {
				return
			}
			index := 0
			shi := &ShareHolerInfo{}
			found := false
			tr.Find(TagNode, "td").Each(func(i int, td *Selection) {
				td.Find(TextNode, "").Each(func(_ int, tn *Selection) {
					if tn.Nodes[0].GetParentNodeTagname() == "td" {
						found = true
						//fmt.Fprintf(os.Stdout, "i-%d, data-%s\n", i, tn.Nodes[0].Root.Data)
						switch index {
						case 0:
							shi.Name = strings.TrimSpace(tn.Nodes[0].Root.Data)
						case 1:
							shi.Count = fmt.Sprintf("%.4f", utility.String2Folat64(strings.TrimSpace(tn.Nodes[0].Root.Data))/10000)
						case 2:
							shi.Ratio = utility.String2Folat32(strings.TrimSpace(tn.Nodes[0].Root.Data))
						}

						index++
					}
				})
			})
			if found {
				//shi.Date = date
				shiList = append(shiList, shi)
			}
		})

	})

	return shiList
}

func (tree *HTMLDoc) JJCC_GetCurrentDate() string {
	var result string

	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {
			return
		}

		table.Find(TagNode, "td").Each(func(i int, td *Selection) {
			td.Find(TagNode, "option").Each(func(i int, option *Selection) {
				value := option.Nodes[0].GetAttrByName("selected")
				if value != "" {
					// Format: <option selected='selected' value='2017-06-30'>2017-06-30</option>
					result = option.Nodes[0].GetAttrByName("value")
				}
			})
		})
	})

	return result
}

func (tree *HTMLDoc) JJCC_ParseRecordsDate() (result []time.Time) { //map[string]bool {
	//result := make(map[string]bool)

	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {
			return
		}

		table.Find(TagNode, "td").Each(func(i int, td *Selection) {

			td.Find(TagNode, "option").Each(func(i int, option *Selection) {
				//result[option.Nodes[0].GetAttrByName("value")] = false
				tmp := option.Nodes[0].GetAttrByName("value")
				if recordDate, err := time.Parse("2006-01-02", tmp); err == nil {
					result = append(result, recordDate)
				} else {
					logger.Warningf("Parse record date failure, %s", tmp)
				}
			})
		})
	})

	return
}
