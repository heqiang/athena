package selector

import (
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
	"strings"
)

type (
	Selector struct {
		Text string
	}

	Node struct {
		node *html.Node
	}
)

func NewSelector(Text string) (*Node, error) {
	node, err := html.Parse(strings.NewReader(Text))
	if err != nil {
		return nil, err
	}

	return &Node{
		node: node,
	}, nil
}

var _ xpath.NodeNavigator = &NodeNavigator{}

func CreateXPathNavigator(top *html.Node) *NodeNavigator {
	return &NodeNavigator{curr: top, root: top, attr: -1}
}
func (s *Node) Xpath(xpExpr string) ([]*Node, error) {
	exp, err := getQuery(xpExpr)
	if err != nil {
		return nil, err
	}

	var selectors []*Node
	t := exp.Select(CreateXPathNavigator(s.node))
	for t.MoveNext() {
		nav := t.Current().(*NodeNavigator)
		n := getCurrentNode(nav)
		selectors = append(selectors, &Node{
			node: n,
		})
	}

	return selectors, nil
}

func (s *Selector) getAllNode() {

}
