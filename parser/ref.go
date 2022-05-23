package parser

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func getRef(c *html.Node, in *Result, l string) {
	var ref struct {
		Type  string
		Value string
		Id    string
	}
	r := len(in.Senses[len(in.Senses)-1].Reference) - 1
	if c.Data == "dl" {
		in.Senses[len(in.Senses)-1].Reference = append(in.Senses[len(in.Senses)-1].Reference, ref)

	} else if CheckClass(c, fmt.Sprintf("manyLang%s", l)) && c.Data == "span" {
		in.Senses[len(in.Senses)-1].Reference[r].Type = strings.TrimSpace(c.FirstChild.Data)

	} else if c.Data == "a" && CheckClass(c, "undL") {
		re := regexp.MustCompile("[0-9]+")
		id := c.Attr[0].Val
		id = re.FindAllString(id, -1)[0]
		in.Senses[len(in.Senses)-1].Reference[r].Id = id
		in.Senses[len(in.Senses)-1].Reference[r].Value = c.FirstChild.Data

	} else if c.Data == "dd" && c.FirstChild.Type == html.TextNode {
		in.Senses[len(in.Senses)-1].Reference[r].Value = strings.TrimSpace(c.FirstChild.Data)

	} else if c.Data == " 자세히 보기 끝 " {
		return
	}

	for e := c.FirstChild; e != nil; e = e.NextSibling {
		// Skip all commment nodes or nodes whose type is "script"
		if e.Type == html.CommentNode || e.Data == "script" {
			continue
		} else {
			getRef(e, in, l)
		}
	}
}
