package htmlparser

import (
//	"golang.org/x/net/html"
)

type ShareHolerInfo struct {
	Name string
	Count string
	Ratio string
}


// Top 10  shareholders
type ShareHolderType uint32
const (
	Major ShareHolderType = iota
	Free
)

func (tree *HTMLDoc) GetShareholder(shType ShareHolderType) []*ShareHolerInfo{
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
							shi.Name = tn.Nodes[0].Root.Data
						case 1:
							shi.Count = tn.Nodes[0].Root.Data
						case 2:
							shi.Ratio = tn.Nodes[0].Root.Data
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


func (tree *HTMLDoc) GetDateList() []string{
	var result []string

	tree.Find(TagNode,"table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {return}

		table.Find(TagNode,"td").Each(func(i int, td *Selection){

			td.Find(TagNode, "option").Each(func(i int, option *Selection){
				//fmt.Println(option.Nodes[0].GetAttrByName("value"))

				result = append(result, option.Nodes[0].GetAttrByName("value"))
			})
		})
	})

	return result
}
