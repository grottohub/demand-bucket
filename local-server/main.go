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
	Info   struct {
		Headers *interface{}
		Body    *interface{}
		Form    *interface{}
		Query   string
	}
}

type HeaderInfo struct {
	Title string
	Desc  string
}

type Page struct {
	Header HeaderInfo
	Bucket []BucketInfo
}

func init() {
	rndr := &render.Renderer{}
	rndr.Init()

	homePage := &Page{
		Header: HeaderInfo{
			Title: "Home",
			Desc:  "DemandBucket is a RequestBin clone built using Golang and Redis.",
		},
	}

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
			rndr.Render(w, "home", homePage)
		} else {
			if q == "inspect" {
				info := cache.GetBucket(path)

				bucketPage := &Page{
					Header: HeaderInfo{
						Title: path,
						Desc:  "DemandBucket is a RequestBin clone built using Golang and Redis.",
					},
					Bucket: []BucketInfo{},
				}

				for _, req := range info {
					parts := strings.Split(req, "\n")

					h := []byte(parts[0])
					b := []byte(parts[1])
					q := parts[3]

					newReq := BucketInfo{}

					json.Unmarshal(h, &newReq.Info.Headers)
					json.Unmarshal(b, &newReq.Info.Body)
					newReq.Info.Query = q

					bucketPage.Bucket = append(bucketPage.Bucket, newReq)
				}

				rndr.Render(w, "bucket", bucketPage)
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
