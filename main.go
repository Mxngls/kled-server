package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/search", search)
	http.HandleFunc("/view", view)

	http.ListenAndServe(":8090", nil)
}
