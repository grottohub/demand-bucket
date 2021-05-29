package localserver

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n > Received request from %v\n", r.Header["X-Forwarded-For"][0])
		fmt.Println(" > need to create new bucket")
		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n > Received request from %v\n", r.Header["X-Forwarded-For"][0])
		path := html.EscapeString(r.URL.Path)
		q := html.EscapeString(r.URL.RawQuery)

		if path == "/" {
			fmt.Println(" > need to render homepage")
		} else {
			if q == "inspect" {
				fmt.Println(" > need to render inspection of bucket", path)
			} else {
				fmt.Println(" > need to add request to bucket", path)
			}
		}

		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})
}

func Start(port int) {
	pStr := fmt.Sprintf(":%v", port)

	fmt.Printf("\n > Starting server and listening on %v...\n\n", pStr)
	log.Fatal(http.ListenAndServe(pStr, nil))
}
