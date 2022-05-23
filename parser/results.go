package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func ParseResult(r io.Reader, l string) (res Search, err error) {
	var s Search
	var i int
	doc, err := html.Parse(r)
	dfsr(doc, &s, &i, l)
	return s, err
}

func dfsr(n *html.Node, in *Search, i *int, l string) {
	if CheckClass(n, "blue ml5") {
		// Get the number of results
		str := GetTextAll(n)
		arr := strings.Split(str, " ")
		for x := 0; x < len(arr); x++ {
			conv, err := strconv.Atoi(arr[0])
			if err != nil {
				continue
			}
			in.ResCount = conv
		}

	} else if n.Data == "dl" {
		var r Result
		in.Results = append(in.Results, r)

	} else if CheckClass(n, "") && n.Data == "dd" {
		in.Results[*i].Inflections = GetContent(n, "sup")

	} else if CheckClass(n, "word_type1_17") {
		// Get the Hangul
		// Get the Id
		in.Results[*i].Hangul = GetContent(n, "sup")
		re := regexp.MustCompile("[0-9]+")
		id := n.Parent.Attr[0].Val
		id = re.FindAllString(id, -1)[0]
		in.Results[*i].Id = id

	} else if n.Data == "span" && n.FirstChild != nil &&
		n.FirstChild.Type == html.TextNode &&
		n.FirstChild.Data[0:1] == "(" {
		// Get the Hanja (if there is one)
		hanja := MatchBetween(n.FirstChild.Data, "(", ")")
		in.Results[*i].Hanja = hanja

	} else if CheckClass(n, "word_att_type1") {
		// Get the Korean word type
		match := MatchBetween(GetTextAll(n.FirstChild), "「", "」")
		in.Results[*i].TypeKr = strings.ToValidUTF8(match, "")

	} else if CheckClass(n, fmt.Sprintf("manyLang%s", l)) &&
		CheckClass(n.Parent, "word_att_type1") {
		// Get the English word type
		match := GetTextSingle(n.FirstChild)
		in.Results[*i].TypeEng = strings.TrimSpace(cleanStringSpecial([]byte(match)))

	} else if CheckClass(n, "search_sub") {
	out:
		// Get the pronounciation and audio file
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode && c.Data[0:1] == "[" {
				in.Results[*i].Pronounciation = c.Data + "]"
			} else if c.Data == "a" {
				for _, a := range c.Attr {
					if a.Key == "href" {
						in.Results[*i].Audio = MatchBetween(a.Val, "'", "')")
						break out
					}
				}
			}
		}
	}

	// Get the level of the word
	if CheckClass(n, "ri-star-s-fill") {
		in.Results[*i].Level++
	}

	// Get the english translation
	// if (CheckClass(n, fmt.Sprintf("manyLang%s mt15", l)) || CheckClass(n, fmt.Sprintf("manyLang%s ", l))) && !CheckClass(n.NextSibling, "") {
	if (CheckClass(n, fmt.Sprintf("manyLang%s ", l))) && GetTextAll(n) != "" {
		s := InitSense()
		in.Results[*i].Senses = append(in.Results[*i].Senses, s)
		l := len(in.Results[*i].Senses)
		in.Results[*i].Senses[l-1].Translation = cleanStringSpecial([]byte(GetTextAll(n)))
	}

	// Get the korean definition
	if CheckClass(n, "ml20") {
		l := len(in.Results[*i].Senses)
		in.Results[*i].Senses[l-1].KrDefinition = GetTextAll(n)
	}

	// Get the english definition
	if CheckClass(n, fmt.Sprintf("manyLang%s ml20", l)) && GetTextAll(n) != "" {
		l := len(in.Results[*i].Senses)
		in.Results[*i].Senses[l-1].Definition = GetTextAll(n)
	}

	// Increment the index by one
	if CheckClass(n, fmt.Sprintf("manyLang%s ml20", l)) && n.NextSibling.NextSibling == nil {
		*i++
	}

	// Get the number of pages
	if CheckClass(n, "paging_area") {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if CheckClass(c, "btn_first") {
				in.Pages = append(in.Pages, -4)
			} else if CheckClass(c, "btn_prev") {
				in.Pages = append(in.Pages, -3)
			} else if CheckClass(c, "paging_num") || CheckClass(c, "paging_num on") {
				page, err := strconv.Atoi(GetTextAll(c))
				if err != nil {
					panic(err)
				}
				in.Pages = append(in.Pages, page)
			} else if CheckClass(c, "btn_next") {
				in.Pages = append(in.Pages, -2)
			} else if CheckClass(c, "btn_last") {
				in.Pages = append(in.Pages, -1)
			}
		}
	}

	// Traverse the tree of nodes vi depth-first search
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Skip all commment nodes or nodes whose type is "script"
		if c.Type == html.CommentNode && c.Data == "script" {
			continue
		} else {
			dfsr(c, in, i, l)
		}
	}
}
