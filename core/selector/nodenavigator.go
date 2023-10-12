package selector

import (
	"fmt"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
	"strings"
)

type NodeNavigator struct {
	root, curr *html.Node
	attr       int
}

func (h *NodeNavigator) NodeType() xpath.NodeType {
	switch h.curr.Type {
	case html.CommentNode:
		return xpath.CommentNode
	case html.TextNode:
		return xpath.TextNode
	case html.DocumentNode:
		return xpath.RootNode
	case html.ElementNode:
		if h.attr != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	case html.DoctypeNode:
		// ignored <!DOCTYPE HTML> declare and as Root-Node type.
		return xpath.RootNode
	}
	panic(fmt.Sprintf("unknown HTML node type: %v", h.curr.Type))
}

func (h *NodeNavigator) LocalName() string {
	if h.attr != -1 {
		return h.curr.Attr[h.attr].Key
	}
	return h.curr.Data
}

func (*NodeNavigator) Prefix() string {
	return ""
}

func (h *NodeNavigator) Value() string {
	switch h.curr.Type {
	case html.CommentNode:
		return h.curr.Data
	case html.ElementNode:
		if h.attr != -1 {
			return h.curr.Attr[h.attr].Val
		}
		return InnerText(h.curr)
	case html.TextNode:
		return h.curr.Data
	}
	return ""
}
func InnerText(n *html.Node) string {
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
	output(&b, n)
	return b.String()
}
func (h *NodeNavigator) Copy() xpath.NodeNavigator {
	n := *h
	return &n
}

func (h *NodeNavigator) MoveToRoot() {
	h.curr = h.root
}

func (h *NodeNavigator) MoveToParent() bool {
	if h.attr != -1 {
		h.attr = -1
		return true
	} else if node := h.curr.Parent; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToNextAttribute() bool {
	if h.attr >= len(h.curr.Attr)-1 {
		return false
	}
	h.attr++
	return true
}

func (h *NodeNavigator) MoveToChild() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.FirstChild; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToFirst() bool {
	if h.attr != -1 || h.curr.PrevSibling == nil {
		return false
	}
	for {
		node := h.curr.PrevSibling
		if node == nil {
			break
		}
		h.curr = node
	}
	return true
}

func (h *NodeNavigator) String() string {
	return h.Value()
}

func (h *NodeNavigator) MoveToNext() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.NextSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToPrevious() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.PrevSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*NodeNavigator)
	if !ok || node.root != h.root {
		return false
	}

	h.curr = node.curr
	h.attr = node.attr
	return true
}

func getCurrentNode(n *NodeNavigator) *html.Node {
	if n.NodeType() == xpath.AttributeNode {
		childNode := &html.Node{
			Type: html.TextNode,
			Data: n.Value(),
		}
		return &html.Node{
			Type:       html.ElementNode,
			Data:       n.LocalName(),
			FirstChild: childNode,
			LastChild:  childNode,
		}

	}
	return n.curr
}
