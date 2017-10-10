package htmlparser

import (
//	"golang.org/x/net/html"

	"fmt"
)

type ShareHolerInfo struct {
	name string
	count int
	ratio int
}


// Top 10 major shareholders
func (tree *HTMLDoc) GetMajorShareholder() {
	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {return}
		if i != 5 {return }  // The major shareholder is in the 5th table.

		table.Find(TagNode,"tr").Each(func(i int, tr *Selection) {
			tr.Find(TagNode, "td").Each(func(i int, td *Selection) {
				fmt.Println(td.Nodes[0].Root.Attr)
			})
		})

	})
}

// Top 10 shareholders
func (tree *HTMLDoc) GetShareholder() {

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
