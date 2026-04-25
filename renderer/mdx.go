package renderer

import (
	"os"
	"text/template"
)

func RenderDoc(doc Doc) {
	tmplFile := "templates/problem.tpl"
	if doc.Type == "concept" {
		tmplFile = "templates/concept.tpl"
	}

	tmpl := template.Must(template.ParseFiles(tmplFile))

	f, _ := os.Create(doc.Output)
	defer f.Close()

	tmpl.Execute(f, doc)
}
