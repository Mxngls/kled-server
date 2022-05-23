package parser

type Search struct {
	Results  []Result
	ResCount int
	Pages    []int
}

type Result struct {
	Id              string
	Hangul          string
	Hanja           string
	TypeKr          string
	TypeEng         string
	Pronounciation  string
	Audio           string
	Level           int
	Inflections     string
	InflectionLinks []inflectionLink
	Senses          []sense
}

type sense struct {
	Translation  string
	Definition   string
	KrDefinition string
	Examples     []string
	Reference    []struct {
		Type  string
		Value string
		Id    string
	}
}

type inflectionLink struct {
	Id     string
	Hangul string
}

func InitSense() sense {
	return sense{}
}

func InitInflectionLinks() inflectionLink {
	return inflectionLink{}
}
