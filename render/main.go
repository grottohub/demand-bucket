package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Renderer struct {
	Template *template.Template
}

func toStrings(arr []interface{}) []string {
	s := []string{}

	for _, v := range arr {
		s = append(s, v.(string))
	}

	return s
}

func (r *Renderer) Init() {
	funcs := template.FuncMap{"join": strings.Join, "toStrings": toStrings}

	var err error
	r.Template, err = template.New("main").Funcs(funcs).ParseGlob("render/templates/*.gohtml")
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) Render(w http.ResponseWriter, n string, data interface{}) {
	fmt.Printf(" > rendering template %v\n", n)

	err := r.Template.ExecuteTemplate(w, n, data)
	if err != nil {
		fmt.Printf(" > encountered error:\n")
		fmt.Fprintf(w, "ERR: %v", err)
		panic(err)
	}
}
