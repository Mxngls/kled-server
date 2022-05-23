package main

import (
	"bytes"
	"io"
	"net/http"

	"dict-wrapper/parser"
)

func search(w http.ResponseWriter, req *http.Request) {

	word := req.FormValue("word")
	lang := req.FormValue("lang")
	langCode := req.FormValue("langCode")
	page := req.FormValue("page")

	resp := requestSearch(word, lang, langCode, page)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(body)
	data, err := parser.ParseResult(reader, langCode)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	JSON, err := JSONMarshal(data)

	w.Write(JSON)
}

func view(w http.ResponseWriter, req *http.Request) {

	id := req.FormValue("id")
	lang := req.FormValue("lang")
	langCode := req.FormValue("langCode")

	resp := requestView(id, lang, langCode)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(body)
	data, err := parser.ParseView(reader, id, langCode)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	JSON, err := JSONMarshal(data)

	w.Write(JSON)
}
