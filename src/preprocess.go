package gojustext

import (
	"html"
	"log"
	"regexp"
	"strings"
)

func preprocess(htmlStr, encoding, defaultEncoding, encErrors string) *html.Node {
	
	root, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	addKwTags(root)
	removeElements(root, []string{"head", "script", "style"})

	return root
}

type nodeIterator func(n *html.Node)

func nodeIter(n *html.Node, f nodeIterator) {
	f(n)
	for _, c := range n.Child {
		nodeIter(c, f)
	}
}

// Bug in here!
// Swaps the order of some elements!
func addKwTags(root *html.Node) {
	var blankText *regexp.Regexp = regexp.MustCompile("^[\n\r\t ]*$")
	var nodesWithText []*html.Node

	var markTextAndTail nodeIterator
	markTextAndTail = func(node *html.Node) {
		if node.Type != html.CommentNode || node.Type != html.DoctypeNode {
			if node.Type == html.TextNode {
				nodesWithText = append(nodesWithText, node)
			}
		}
	}
	nodeIter(root, markTextAndTail)

	for _, node := range nodesWithText {
		if blankText.MatchString(node.Data) {
			node.Data = ""
		} else {
			kw := &html.Node{
				Parent: nil,
				Type: html.ElementNode,
				Data: "kw",
			}
			kw.Child = append(kw.Child, node)
			replaceNode(node, kw)
		}
	}
}

func removeElements(root *html.Node, elementsToRemove []string) {
	var toBeRemoved []*html.Node
	var markRemovableNodes = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, nodeName := range elementsToRemove {
				if node.Data == nodeName {
					toBeRemoved = append(toBeRemoved, node)
				}
			}
		}
	}
	nodeIter(root, markRemovableNodes)

	for _, node := range toBeRemoved {
		node.Parent.Remove(node)
	}
}

// replaceNode replaces a Node in a Node tree with a replacement Node.
// Should be moved into html/utils package
func replaceNode(originalNode *html.Node, newNode *html.Node) {
	slice := originalNode.Parent.Child
	for position, n := range slice {
		if n == originalNode {
			slice = append(slice[:position], append([]*html.Node{newNode}, slice[position:]...)...)
			return
		}
	}
}
