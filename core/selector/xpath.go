package selector

import (
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
	"io"
)

type (
	Selector struct {
		node *html.Node
	}
)

func NewSelector(r io.Reader) (*Selector, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &Selector{
		node: node,
	}, nil
}

var _ xpath.NodeNavigator = &NodeNavigator{}

func CreateXPathNavigator(top *html.Node) *NodeNavigator {
	return &NodeNavigator{curr: top, root: top, attr: -1}
}
func (s *Selector) Xpath(xpExpr string) ([]*Selector, error) {
	exp, err := getQuery(xpExpr)
	if err != nil {
		return nil, err
	}

	var selectors []*Selector
	t := exp.Select(CreateXPathNavigator(s.node))
	for t.MoveNext() {
		nav := t.Current().(*NodeNavigator)
		n := getCurrentNode(nav)
		selectors = append(selectors, &Selector{
			node: n,
		})
	}

	return selectors, nil
}

func (s *Selector) getAllNode() {

}
