package parser

import (
	"strings"

	"golang.org/x/net/html"
)

// Check if a given attribute exists
func CheckAttribute(n *html.Node, key string) bool {
	for _, el := range n.Attr {
		if el.Key == key {
			return true
		}
	}
	return false
}

// Check the value of a given attribute
func CheckAttributeVal(n *html.Node, key string, val string) bool {
	for _, el := range n.Attr {
		if el.Key == key {
			if el.Val == val {
				return true
			}
		}
	}
	return false
}

// Check for the existence of a given id
func CheckId(n *html.Node, id string) bool {
	if CheckAttributeVal(n, "id", id) {
		return true
	}
	return false
}

// Check for the existence of a given class
func CheckClass(n *html.Node, class string) bool {
	if CheckAttributeVal(n, "class", class) {
		return true
	}
	return false
}

// Get all contents of a given node
func GetContent(n *html.Node, tag string) (content string) {
	if n.Data == tag {
		return "<" + n.Data + ">" + n.FirstChild.Data + "</" + tag + ">"
	} else if n.Data == "br" {
		return "<" + "br" + ">"
	} else if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content += GetContent(c, tag)
	}
	return strings.Join(strings.Fields(content), " ")
}

// Get the text contents of all text nodes of a given node
func GetTextAll(n *html.Node) (text string) {
	if n.Type == html.TextNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += GetTextAll(c) + " "
	}
	return strings.Join(strings.Fields(text), " ")
}

// Get the text contents of a single given node
func GetTextSingle(n *html.Node) (text string) {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return strings.TrimSpace(c.Data)
		}
	}
	return ""
}

// Get the sub-string of a given string that is marked by a specified
// open and close character
func MatchBetween(s string, open string, close string) string {
	i := strings.Index(s, open)
	if i >= 0 {
		j := strings.Index(s, close)
		if j >= 0 {
			return s[i+1 : j]
		}
	}
	return ""
}

// Get rid of everything other than alphanumeric characters in
// given string
func cleanString(s []byte) string {
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func cleanStringSpecial(s []byte) string {
	j := 0
	for _, b := range s {
		if (b != '.') && !('1' <= b && b <= '9') {
			s[j] = b
			j++
		}
	}
	return strings.TrimSpace(string(s[:j]))
}

// Get rid of redundant spaces in a given string
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), "")
}
