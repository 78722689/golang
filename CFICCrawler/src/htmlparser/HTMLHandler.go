package htmlparser

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"regexp"
	//"go/doc"
)

type FindType uint32
const (
	TagNode FindType = iota
	TextNode
	Attr
)

type StockInfo struct {
	Number string
	Name string
	Link string
}

type Selection struct {
	//Nodes []*html.Node
	Nodes []*HTMLDoc
	PreSel *Selection
}

type HTMLDoc struct {
	Root *html.Node
	*Selection
	//stockList []StockInfo
}

func (doc *HTMLDoc)GetAttrByName(name string) string{
	if len(doc.Root.Attr) == 0 {return ""}

	for _, value := range doc.Root.Attr {
		if value.Key == name {return  value.Val}
	}

	return ""

}

func (tree *HTMLDoc) isStockName(value string) []byte {
	reg, _ := regexp.Compile(`\([\d]{6}\)$`)

	return reg.Find([]byte(value))
}

func (tree *HTMLDoc) getStockNameByValue(value string) []byte {
	reg, _ := regexp.Compile(`[A-Z]*[\*]*[A-Z]*[\p{Han}]+`)

	return reg.Find([]byte(value))
}

func (tree *HTMLDoc) getStockNumberByValue(value string) []byte {
	reg, _ := regexp.Compile(`[\d]{6}`)

	return reg.Find([]byte(value))
}

func (tree *HTMLDoc) GetAllStocks() []StockInfo {
	var loopnode func(*html.Node)
	var stockList []StockInfo

	// Loop all nodes to lookup the node where the name matched.
	loopnode = func(node *html.Node) {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if  len(tree.isStockName(child.Data)) > 0 &&
				child.Parent.Data == "a" &&
				len(child.Parent.Attr) > 0 &&
				child.Parent.Attr[0].Key == "href" {

				// Put stock details to array.
				stockList = append(stockList, StockInfo{
						Link:child.Parent.Attr[0].Val,
						Number:string(tree.getStockNumberByValue(child.Data)),
						Name:string(tree.getStockNameByValue(child.Data)),
					})


					fmt.Fprintf(os.Stdout, "Found link:%s len-%d number:%s name:%s\r\n",
						child.Parent.Attr[0].Val,
						len(child.Data),
						tree.getStockNumberByValue(child.Data),
						tree.getStockNameByValue(child.Data))

			}

			loopnode(child)
		}
	}

	loopnode(tree.Root)

	return stockList
}


func ParseFromFile(file string) (*HTMLDoc, error) {

	tree := &HTMLDoc{}
	contents,err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read file %v failed\n", file)
		return tree, err
	}

	doc,err := html.Parse(strings.NewReader(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse file %v failed\n", file)
		return tree, err
	}
	tree.Root = doc
	tree.Selection = &Selection{[]*HTMLDoc{tree},
								nil}
	return tree, nil
}

func ParseFromNode(root *html.Node) (*HTMLDoc, error) {
	tree := &HTMLDoc{}
	tree.Root = root

	return tree, nil
}


func loopNode(node *HTMLDoc, tag string) []*HTMLDoc {
	var result []*HTMLDoc

	for child := node.Root.FirstChild; child != nil; child = child.NextSibling {
		doc := &HTMLDoc{child, nil}
		if child.Type == html.ElementNode && child.Data == tag {
			result = append(result, doc)
		}

		result = append(result, loopNode(doc, tag)...)
	}

	return result
}

func findByTag(node *HTMLDoc, tag string) []*HTMLDoc {
	var result []*HTMLDoc

	for child := node.Root.FirstChild; child != nil; child = child.NextSibling {
		doc := &HTMLDoc{child, nil}
		if child.Type == html.ElementNode && child.Data == tag {
			result = append(result, doc)
		}

		result = append(result, loopNode(doc, tag)...)
	}

	return result
}

func findByText(node *HTMLDoc, filter string) []*HTMLDoc {

	return nil
}

func (sel *Selection)Find(mode FindType, filter string) *Selection {
	var nodes []*HTMLDoc

	for _, doc := range sel.Nodes {
		switch mode {
		case TagNode:
			nodes = append(nodes, findByTag(doc, filter)...)
		case TextNode:
			nodes = append(nodes, findByText(doc, filter)...)
		}
	}


	return &Selection{nodes, sel}
}

func (sel *Selection)Each(f func(int, *Selection)) *Selection {
	for i, node := range sel.Nodes {
		s := &Selection{[]*HTMLDoc{node}, sel}
		f(i, s)
	}

	return sel
}

func (sel *Selection)GetNodeByAttr(name string, value string) []*HTMLDoc {
	var result []*HTMLDoc

	for _, n := range sel.Nodes {
		for _, a := range n.Root.Attr {
			if a.Key == name && a.Val == value {
				result = append(result, n)
			}
		}
	}

	return result
}