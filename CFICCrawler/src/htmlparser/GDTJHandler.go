package htmlparser

import (
	"golang.org/x/net/html"
	"fmt"
	//"os"
)

func (tree *HTMLDoc) Get() {
	var loopnode func(*html.Node)

	loopnode = func(node *html.Node) {
		for child := node.FirstChild; child != nil; child = child.NextSibling {

			//if child.Type == html.ElementNode {
				fmt.Println(child.Namespace)
			//}

			loopnode(child)
		}

	}

	loopnode(tree.Doc)
}
