package localserver

import (
	"demand-bucket/cache"
	"demand-bucket/render"
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

type BucketInfo struct {
	Bucket string
	Info   *interface{}
}

func init() {
	rndr := &render.Renderer{}
	rndr.Init()

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n > Received request from %v\n", r.Header["X-Forwarded-For"][0])
		cache.AddBucket()
		res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
		fmt.Fprint(w, res)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n > Received request from %v\n", r.Header["X-Forwarded-For"][0])
		path := html.EscapeString(r.URL.Path)
		q := html.EscapeString(r.URL.RawQuery)

		if path == "/" {
			rndr.Render(w, "home", nil)
		} else {
			if q == "inspect" {
				info := cache.GetBucket(path)
				b := &BucketInfo{Bucket: path}

				req := strings.Split(info[0], "\n")
				err := json.Unmarshal([]byte(req[0]), &b.Info)
				if err != nil {
					panic(err)
				}

				fmt.Printf(" > unmarshaled: %+v\n", b)

				rndr.Render(w, "bucket", b)
			} else {
				req := formatRequest(r)
				cache.AddRequest(path, req)
				res := fmt.Sprintf("ip:%v\n", r.Header["X-Forwarded-For"][0])
				fmt.Fprint(w, res)
			}
		}
	})
}

func Start(port int) {
	pStr := fmt.Sprintf(":%v", port)

	fmt.Printf(" > Starting server and listening on %v...\n\n", pStr)
	log.Fatal(http.ListenAndServe(pStr, nil))
}
