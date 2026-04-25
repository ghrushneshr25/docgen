package renderer

import (
	"os"
	"text/template"
)

type IndexItem struct {
	Title      string
	Link       string
	Difficulty string
	Tags       string
}

type Index struct {
	Title         string
	Slug          string
	ConceptItems  []IndexItem
	ProblemItems  []IndexItem
	Path          string
}

func RenderIndex(idx Index) {
	t := template.Must(template.ParseFiles("templates/index.tpl"))
	f, _ := os.Create(idx.Path)
	defer f.Close()
	t.Execute(f, idx)
}
