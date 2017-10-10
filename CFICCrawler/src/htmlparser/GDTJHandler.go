package htmlparser

import (
	"golang.org/x/net/html"
	"fmt"

)

type ShareHolerInfo struct {
	name string
	count int
	ratio int
}


// Top 10 major shareholders
/*func (tree *HTMLDoc) GetMajorShareholder() {
	tree.Find(TagName, "table").Each(func(i int, s *Selection) {
		if i != 5 && GetAttrByName(s.Doc.Root, "id") != "tabh" {return }
		sel := s.Find(TagName,"td")
		for _, n := range sel.Nodes {
			fmt.Println(n.Data)
		}
	})
}
*/
// Top 10 shareholders
func (tree *HTMLDoc) GetShareholder() {

}

func (tree *HTMLDoc) GetDateList() []string{
	var result []string

	s := tree.Find(TagName,"table").Each(func(i int, s *Selection) {
		if GetAttrByName(s.Doc.Root, "id") != "tabh" {return }

		//fmt.Fprintf(os.Stdout, "table:%d  attr:%v  element-len:%d\n", i, s.Doc.Root.Attr, len(s.Nodes))
		s.Find(TagName,"td").Each("option", func(i int, s *Selection){
			/*fmt.Fprintf(os.Stdout, "td:%d  attr:%v  element-len:%d\n",
				i, :=
				s.Doc.Root.Attr,
				len(s.Nodes))
			*/

			if len(s.Nodes) > 0 {
				for _, n := range s.Nodes {

					result= append(result, GetAttrByName(n, "value"))
					//fmt.Println(n.Attr)
				}
			}
		})
	})

	return result
}
