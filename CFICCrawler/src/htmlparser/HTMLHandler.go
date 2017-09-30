package htmlparser

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"regexp"
)

type StockInfo struct {
	Number string
	Name string
	Link string
}

type Sellector struct {
	Nodes []*html.Node
}

type HTMLDoc struct {
	Root *html.Node
	*Sellector
	//stockList []StockInfo
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
	tree.Sellector = &Sellector{[]*html.Node{doc}}
	return tree, nil
}

func ParseFromNode(root *html.Node) (*HTMLDoc, error) {
	tree := &HTMLDoc{}
	tree.Root = root

	return tree, nil
}

func findByTag(node *html.Node, tag string) []*html.Node {
	var result []*html.Node

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			result = append(result, child)
		}

		result = append(result, loopNode(child, tag)...)
	}

	return result
}

func (sel *Sellector)Find(node *html.Node, tag string) *Sellector {
	return &Sellector{findByTag(node, tag)}
}

func (sel *Sellector)Each(tag string, f func(int, *Sellector)) *Sellector {
	for i, node := range sel.Nodes {
		fmt.Println(node.Attr)
		nodes := findByTag(node, tag)
		f(i, &Sellector{nodes})
	}

	return sel
}