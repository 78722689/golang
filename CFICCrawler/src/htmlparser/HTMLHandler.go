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

type HTMLDoc struct {
	Doc *html.Node
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

	loopnode(tree.Doc)

	return stockList
}


func ParseFromFile(file string) (*HTMLDoc, error) {

	tree := &HTMLDoc{}
	contents,err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read file %v failed", file)
		return tree, err
	}

	doc,err := html.Parse(strings.NewReader(string(contents)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse file %v failed", file)
		return tree, err
	}
	tree.Doc = doc

	return tree, nil
}

func ParseFromNode(root *html.Node) (*HTMLDoc, error) {
	tree := &HTMLDoc{}
	tree.Doc = root

	return tree, nil
}