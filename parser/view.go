package parser

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func ParseView(result_html io.Reader, id string, l string) (res Result, err error) {
	var r Result
	r.Id = fmt.Sprintf("%s", id)
	doc, err := html.Parse(result_html)
	dfsv(doc, &r, l)
	return r, err
}

func dfsv(n *html.Node, in *Result, l string) *html.Node {
	// Get the Hangul and Hanja
	if CheckClass(n, "word_head") {
		in.Hangul = GetContent(n, "sup")
		hanja := MatchBetween(n.NextSibling.Data, "(", ")")
		in.Hanja = hanja

	} else if CheckClass(n, "search_sub") {
		// Get the pronounciation and audio file
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode && strings.TrimSpace(c.Data) != "" && c.Data[0:1] == "[" {
				in.Pronounciation = strings.TrimSpace(c.Data) + "]"
			} else if c.Data == "a" && c.Attr[1].Val == "sound" {
				for _, a := range c.Attr {
					if a.Key == "href" {
						in.Audio = MatchBetween(a.Val, "'", "');")
						break
					}
				}
				break
			}
		}

	} else if CheckClass(n, "word_att_type1") {
		// Get the Korean word type
		match := MatchBetween(GetTextAll(n.FirstChild), "「", "」")
		in.TypeKr = strings.ToValidUTF8(match, "")

	} else if CheckClass(n, fmt.Sprintf("manyLang%s ml5", l)) {
		// Get the English word type
		in.TypeEng = strings.TrimSpace(GetTextSingle(n.FirstChild))

	} else if CheckClass(n, "ri-star-s-fill") {
		// Get the level of the word
		in.Level++

	} else if CheckClass(n, "printArea") {
		// Get the Inflections
		in.Inflections = (GetTextSingle(n))
		inf := InitInflectionLinks()
		ind := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "a" {
				in.InflectionLinks = append(in.InflectionLinks, inf)
				re := regexp.MustCompile("[0-9]+")
				id := c.Attr[1].Val
				id = re.FindAllString(id, -1)[0]
				in.InflectionLinks[ind].Id = id
				in.InflectionLinks[ind].Hangul = GetTextSingle(c) + "<sup>" + GetTextSingle(c.NextSibling) + "</sup>"
				ind++
			}
		}

	} else if CheckClass(n, fmt.Sprintf("multiTrans manyLang%s mb10 printArea", l)) {
		// Get the english translation
		s := InitSense()
		in.Senses = append(in.Senses, s)
		l := len(in.Senses)
		in.Senses[l-1].Translation = cleanStringSpecial([]byte(GetTextAll(n)))

	} else if CheckClass(n, "senseDef ml20") {
		// Get the korean definition
		l := len(in.Senses)
		in.Senses[l-1].KrDefinition = GetTextAll(n)

	} else if CheckClass(n, fmt.Sprintf("multiSenseDef manyLang%s ml20", l)) {
		// Get the english definition
		in.Senses[len(in.Senses)-1].Definition = GetTextAll(n)
		if n.NextSibling.NextSibling == nil {
		}

	} else if CheckClass(n, "dot") || CheckClass(n, "dot printArea") {
		// Get the examples
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "li" {
				ex := GetContent(c, "b")
				in.Senses[len(in.Senses)-1].Examples = append(in.Senses[len(in.Senses)-1].Examples, ex)
			} else if c.Parent.Data != "ul" && c.Parent.Data != "li" {
				break
			}
		}

	} else if CheckClass(n, "heading_wrap") {
		getRef(n, in, l)
	}

	// Traverse the tree of nodes vi depth-first search
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Skip all commment nodes or nodes whose type is "script"
		if c.Type == html.CommentNode || c.Data == "script" {
			continue
		} else {
			dfsv(c, in, l)
		}
	}
	return n
}
