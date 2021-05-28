package localserver

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func init() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request:\n\n")
		prettyPrint(r.Header)

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
}

func Start() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Status() {
	fmt.Println("TODO: server status details")
}
