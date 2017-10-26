package htmlparser

import (
	"httpcontroller"
	"strings"
	"utility"
)

type ShareHolerInfo struct {
	Name string
	Count int64
	Ratio float32
}


// Top 10  shareholders
type ShareHolderType uint32
const (
	Major ShareHolderType = iota
	Free
)

func (tree *HTMLDoc)GDTJ_Request(url string, file string) (*HTMLDoc, error){
	request := httpcontroller.Request{
		//Proxy:&httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"},
		Url : url,
		File : file,
		OverWrite : false,
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

func (tree *HTMLDoc) GDTJ_GetShareholder(shType ShareHolderType) []*ShareHolerInfo{
	var shiList []*ShareHolerInfo

	// Default to read table 6 for common shareholders.
	tableID := 6
	if shType == Major {tableID = 5}

	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {return}
		if i != tableID {return }  // The major shareholder is in the 5th table.

		table.Find(TagNode,"tr").Each(func(i int, tr *Selection) {
			if i<=1 {return }
			index := 0
			shi := &ShareHolerInfo{}
			found := false
			tr.Find(TagNode, "td").Each(func(i int, td *Selection) {
				td.Find(TextNode, "").Each(func(_ int, tn *Selection) {
					if tn.Nodes[0].GetParentNodeTagname() == "td"{
						found = true
						//fmt.Fprintf(os.Stdout, "i-%d, data-%s\n", i, tn.Nodes[0].Root.Data)
						switch index {
						case 0:
							shi.Name = strings.TrimSpace(tn.Nodes[0].Root.Data)
						case 1:
							shi.Count = utility.String2Int64(strings.TrimSpace(tn.Nodes[0].Root.Data))
						case 2:
							shi.Ratio = utility.String2Folat32(strings.TrimSpace(tn.Nodes[0].Root.Data))
						}

						index ++
					}
				})
			})
			if found {
				shiList = append(shiList, shi)
			}
		})

	})

	return shiList
}

func (tree *HTMLDoc) GetCurrentDate() string{
	var result string

	tree.Find(TagNode,"table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {return}

		table.Find(TagNode,"td").Each(func(i int, td *Selection){
			td.Find(TagNode, "option").Each(func(i int, option *Selection){
				value := option.Nodes[0].GetAttrByName("selected")
				if value != "" {
					// Format: <option selected='selected' value='2017-06-30'>2017-06-30</option>
					result =  option.Nodes[0].GetAttrByName("value")
				}
			})
		})
	})

	return result
}


func (tree *HTMLDoc) GetDateList() map[string]bool{
	result := make(map[string]bool)

	tree.Find(TagNode,"table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {return}

		table.Find(TagNode,"td").Each(func(i int, td *Selection){

			td.Find(TagNode, "option").Each(func(i int, option *Selection){
				result[option.Nodes[0].GetAttrByName("value")] = false
			})
		})
	})

	return result
}
