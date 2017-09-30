package htmlparser

import (
	"golang.org/x/net/html"
	"fmt"
	//"os"
	"os"
)

func loopNode(node *html.Node, tag string) []*html.Node {
	var result []*html.Node

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			result = append(result, child)
		}

		result = append(result, loopNode(child, tag)...)
	}

	return result
}

func (tree *HTMLDoc) Get() {
	tree.Find(tree.Root, "table").Each("tr", func(i int, sellector *Sellector) {
		fmt.Fprintf(os.Stdout, "table:%d  element-len:%d\n", i, len(sellector.Nodes))
	})

/*
	for _, table := range result.Nodes {
		fmt.Fprintf(os.Stdout, "%v\n", table)
	}
*/
}
