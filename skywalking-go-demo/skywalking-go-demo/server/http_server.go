package main

import (
	"fmt"
	_ "github.com/apache/skywalking-go"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("Server failed to start: ", err)
	}
}
