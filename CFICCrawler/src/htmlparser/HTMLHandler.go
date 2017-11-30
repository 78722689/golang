package htmlparser

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"regexp"
	"utility"
)
var logger = utility.GetLogger()

type FindType uint32
const (
	TagNode FindType = iota
	TextNode
	Attr
	ErrNode
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

func (doc *HTMLDoc) HaveChildNode() bool {
	if doc.Root.FirstChild == nil {
		return false
	}

	return true
}

func (doc *HTMLDoc) getFirstChildNodeType() FindType {
	switch doc.Root.FirstChild.Type {
	case html.TextNode:
		return TextNode
	case html.ElementNode:
		return TagNode
	}

	return ErrNode
}

func (doc *HTMLDoc)GetData() string {
	return doc.Root.Data
}


func (doc *HTMLDoc)GetParentNodeTagname() string {
	return doc.Root.Parent.Data
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

// Get the stocks info by filter. If filter is empty, get all stocks info
func (tree *HTMLDoc) GetStocks(filterIDs []string) []StockInfo {
	var loopnode func(*html.Node)
	var stockList []StockInfo

	// Loop all nodes to lookup the node where the name matched.
	loopnode = func(node *html.Node) {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if  len(tree.isStockName(child.Data)) > 0 &&
				child.Parent.Data == "a" &&
				len(child.Parent.Attr) > 0 &&
				child.Parent.Attr[0].Key == "href" {

					stockID := string(tree.getStockNumberByValue(child.Data))
					if stockID != "" {
						si := StockInfo{
							Link:   child.Parent.Attr[0].Val,
							Number: stockID,
							Name:   string(tree.getStockNameByValue(child.Data)),
						}

						if len(filterIDs) == 0 {
							stockList = append(stockList, si)
						} else if utility.Contains(filterIDs, stockID) {
							logger.DEBUG("found........................................")
							stockList = append(stockList, si)
						}

						logger.DEBUG(fmt.Sprintf("Found link:%s len-%d number:%s name:%s", si.Link, len(child.Data), si.Number, si.Name))
					} else {
						logger.WARN(fmt.Sprintf("Not found stock ID in node, %v", child))
					}
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


func loopNodeByTag(node *HTMLDoc, nodeType html.NodeType,filter string) []*HTMLDoc {
	var result []*HTMLDoc

	for child := node.Root.FirstChild; child != nil; child = child.NextSibling {
		doc := &HTMLDoc{child, nil}
		if child.Type == nodeType && child.Data == filter {
			result = append(result, doc)
		}

		result = append(result, loopNodeByTag(doc, nodeType, filter)...)
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

		result = append(result, loopNodeByTag(doc, html.ElementNode, tag)...)
	}

	return result
}

func findByText(node *HTMLDoc, filter string) []*HTMLDoc {
	var result []*HTMLDoc

	var loopNode func(*HTMLDoc, html.NodeType, string) []*HTMLDoc

	loopNode = func(node *HTMLDoc, nodeType html.NodeType, filter string) []*HTMLDoc {
		var r []*HTMLDoc

		for child := node.Root.FirstChild; child != nil; child = child.NextSibling {
			doc := &HTMLDoc{child, nil}
			if child.Type == nodeType && strings.TrimSpace(child.Data) != "" {
				if filter != "" && strings.Contains(child.Data, filter) {
					r = append(r, doc)
				} else {
					r = append(r, doc)
				}
			}

			r = append(r, loopNode(doc, nodeType, filter)...)
		}

		return r
	}

	result = append(result, loopNode(node, html.TextNode, filter)...)


	return result
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
