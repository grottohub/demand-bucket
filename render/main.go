package render

import (
	"fmt"
	"html/template"
	"net/http"
)

type Renderer struct {
	Templates map[string]*template.Template
}

func (r *Renderer) Init() {
	r.Templates = map[string]*template.Template{}

	tmpl, err := template.New("home.gohtml").ParseFiles("render/templates/home.gohtml")

	if err != nil {
		panic(err)
	}

	r.Add(tmpl, "home")

	tmpl, err = template.New("bucket.gohtml").ParseFiles("render/templates/bucket.gohtml")

	if err != nil {
		panic(err)
	}

	r.Add(tmpl, "bucket")
}

func (r *Renderer) Add(t *template.Template, n string) {
	r.Templates[n] = t
}

func (r *Renderer) Render(w http.ResponseWriter, n string, data interface{}) {
	t, e := r.Templates[n]

	fmt.Printf(" > rendering template %v\n", t.Name())

	if e {
		err := t.Execute(w, data)
		if err != nil {
			fmt.Printf(" > encountered error:\n")
			fmt.Fprintf(w, "ERR: %v", err)
			panic(err)
		}
	}
}
