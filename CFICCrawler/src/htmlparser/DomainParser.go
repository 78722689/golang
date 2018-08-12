package htmlparser

import (
	"strings"
)

func (tree *HTMLDoc) GetDomains() []string {

	var result []string

	tree.Find(TagNode, "table").Each(func(i int, table *Selection) {
		if len(table.GetNodeByAttr("id", "tabh")) == 0 {
			return
		}

		table.Find(TagNode, "tr").Each(func(i int, tr *Selection) {
			if i <= 2 {
				return
			}

			tr.Find(TagNode, "td").Each(func(j int, td *Selection) {
				// The first column
				if j == 0 {
					td.Find(TextNode, "").Each(func(_ int, tn *Selection) {
						result = append(result, strings.TrimSpace(tn.Nodes[0].Root.Data))
						return
					})
				}
			})
		})

	})

	return result
}
