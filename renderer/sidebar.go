package renderer

import (
	"os"
	"text/template"
)

type SidebarCategory struct {
	Name     string // display label (title-cased folder name)
	Slug     string // docs folder name, e.g. linkedlist → doc id linkedlist/index
	Concepts []string
	Problems []string
}

func RenderSidebar(data []SidebarCategory, out string) {
	t := template.Must(template.ParseFiles("templates/sidebar.tpl"))
	f, _ := os.Create(out)
	defer f.Close()
	t.Execute(f, data)
}
