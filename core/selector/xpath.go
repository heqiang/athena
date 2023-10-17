package selector

import (
	"fmt"
	"github.com/antchfx/xpath"
	"github.com/zeromicro/go-zero/core/logx"
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

// Xpath get all h5 label
func (s *Node) Xpath(xpExpr string) []*Node {
	exp, err := getQuery(xpExpr)
	if err != nil {
		logx.Errorf("xpath error:%v", err)
		return []*Node{}
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
	if len(selectors) == 0 {
		logx.Errorf(fmt.Sprintf(`current xpath:%s get empty node,please check you xpath`, xpExpr))
	}

	return selectors
}

// Text Obtain all text under the current h5 label
func (s *Node) Text() string {
	var output func(*strings.Builder, *html.Node)
	output = func(b *strings.Builder, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			b.WriteString(n.Data)
			return
		case html.CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(b, child)
		}
	}

	var b strings.Builder
	output(&b, s.node)
	return strings.TrimSpace(b.String())
}

// FirstNode get the first node
func (s *Node) FirstNode(xpExpr string) *Node {
	nodes := s.Xpath(xpExpr)
	if len(nodes) == 0 {
		return &Node{}
	}

	return nodes[0]
}

// GetAttribute get the h5 tag attribute value
func (s *Node) GetAttribute(attrName string) string {
	if s.node == nil {
		return ""
	}

	if s.node.Type == html.ElementNode && s.node.Parent == nil && attrName == s.node.Data {
		return InnerText(s.node)
	}

	for _, attr := range s.node.Attr {
		if attr.Key == attrName {
			return strings.TrimSpace(attr.Val)
		}
	}

	return "no such attribute"

}
