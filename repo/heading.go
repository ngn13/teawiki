package repo

import (
	"strings"

	"golang.org/x/net/html"
)

type Heading struct {
	Name     string
	ID       string
	Depth    int
	Parent   *Heading
	Children []*Heading
}

func (h *Heading) Add(child *Heading) {
	if nil == h.Children {
		h.Children = []*Heading{}
	}

	h.Children = append(h.Children, child)
}

var headings = []string{
	"h1", "h2", "h3", "h4", "h5",
}

func findAttr(node *html.Node, name string) string {
	for i := range node.Attr {
		if node.Attr[i].Key == name {
			return node.Attr[i].Val
		}
	}

	return ""
}

func parseNode(node *html.Node, root *[]*Heading, cur *Heading) *Heading {
	var (
		depth int    = -1
		id    string = findAttr(node, "id")
	)

	if node.Type != html.ElementNode || node.FirstChild == nil || id == "" {
		return cur
	}

	for i, h := range headings {
		if node.Data == h {
			depth = i
		}
	}

	if depth == -1 {
		return cur
	}

	heading := &Heading{
		Name:     node.FirstChild.Data,
		ID:       id,
		Depth:    depth,
		Parent:   nil,
		Children: nil,
	}

	for cur != nil && cur.Depth >= heading.Depth {
		cur = cur.Parent
	}

	if cur == nil {
		*root = append(*root, heading)
		return heading
	}

	heading.Parent = cur
	cur.Add(heading)
	return cur
}

func Headings(content string) []*Heading {
	var (
		cur  *Heading = nil
		root          = []*Heading{}
	)

	// document
	node, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return root
	}

	// <html>
	if node = node.FirstChild; node == nil {
		return root
	}

	// <body>
	if node = node.FirstChild.NextSibling; node == nil {
		return root
	}

	for node = node.FirstChild; node != nil; node = node.NextSibling {
		cur = parseNode(node, &root, cur)
	}

	return root
}
