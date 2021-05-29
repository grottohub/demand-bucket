package localserver

import (
	"demand-bucket/cache"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
)

func formatRequest(r *http.Request) string {
	w := strings.Builder{}

	j := json.NewEncoder(&w)

	j.SetEscapeHTML(false)

	j.Encode(r.Header)
	j.Encode(r.Body)
	j.Encode(r.Form)
	j.Encode(r.URL.RawQuery)

	return w.String()
}

func init() {
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n > Received request from %v\n", r.Header["X-Forwarded-For"][0])
		// fmt.Println(" > need to create new bucket")
		cache.AddBucket()
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
				fmt.Printf(" > currently in %v: %v", path, cache.GetBucket(path))
				fmt.Println(" > need to render inspection of bucket", path)
			} else {
				req := formatRequest(r)
				cache.AddRequest(path, req)
				// fmt.Println(" > need to add request to bucket", path)
			}
		}

		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})
}

func Start(port int) {
	pStr := fmt.Sprintf(":%v", port)

	fmt.Printf(" > Starting server and listening on %v...\n\n", pStr)
	log.Fatal(http.ListenAndServe(pStr, nil))
}
