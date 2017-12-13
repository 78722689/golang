package htmlparser

import (
	"golang.org/x/net/html"
	//"fmt"
	"strings"
	"regexp"
)

// Return the modules in the stock home page
func GetStockModules() map[string]int{
	// All modules in stock home page.
	return map[string]int {
		"qgqp": 1,
		"yjyg": 2,
		"yjkb": 3,
		"fxysg": 4,
		"same_hy": 5,
		"same_dq":6,
		"same_gn": 7,
		"ybyl": 8,
		"zcfzb_x": 9,
		"lrfpb_x": 10,
		"xjll": 11,
		"FJCXSY_HJ": 12,
		"gsda": 13,
		"zyfb": 14,
		"ggyl": 15,
		"zlsj": 16,
		"dzjy": 17,
		"yxsb": 18,
		"rzrq": 19,
		"gbjg": 20,
		"gdtj": 21,
		"gdhs": 22,
		"jjcc": 23,
		"jjccbd": 24,
		"fhpx": 25,
		"pgzf": 26,
		"jyzj": 27,
		"tzzk": 28,
		"cbgg": 29,
		"ggqw": 30,
		"zdsj": 31,
		"cwzbMgfxzb": 32,
		"cwzbHlnl": 33,
		"cwzbCznl": 34,
		"cwzbFznl": 35,
		"cwzbJynl": 36,
		"cwzbXjllfx": 37,
		"cwzbFhnl": 38,
		"cwzbZbjg": 39,
		"cwzbLrgc": 40,
		"cwzbDbfx": 41,
	}
}

func (doc *HTMLDoc)GetStockId() string {
	var loopnode func(node *html.Node) string
	loopnode = func(node *html.Node) string {
		var result string
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "a"{
				for _, attr := range child.Attr {
					if attr.Key == "href"{
						if strings.Contains(attr.Val, "stockid=") {
							reg, _ := regexp.Compile(`[\d]{1,6}`)
							return string(reg.Find([]byte(attr.Val)))
						}
					}
				}
			}

			result = loopnode(child)
			if result != "" {return result}
		}

		return ""
	}

	return loopnode(doc.Root)
}

// Get all modules URL
func (doc *HTMLDoc)GetModuleURL(filter string) []string {
	var loopnode func(node *html.Node)
	var urls []string

	loopnode = func(node *html.Node) {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Data == "a" &&
				len(child.Attr) >= 2 &&
				child.Attr[1].Key == "href" &&
				strings.Contains(child.Attr[1].Val, filter) {
					urls = append(urls, child.Attr[1].Val)
			}

			loopnode(child)
		}
	}
	loopnode(doc.Root)

	return urls
}