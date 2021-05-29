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
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("need to create new bucket")
		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request:\n\n")
		prettyPrint(r.Header)
		path := html.EscapeString(r.URL.Path)
		q := html.EscapeString(r.URL.RawQuery)

		if path == "/" {
			fmt.Println("need to render homepage")
		} else {
			if q == "inspect" {
				fmt.Println("need to render inspection of bucket", path)
			} else {
				fmt.Println("need to add request to bucket", path)
			}
		}

		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})
}

func Start() {
	fmt.Println("\n > Starting server and listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Status() {
	fmt.Println("TODO: server status details")
}
