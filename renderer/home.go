package renderer

import (
	"os"
	"text/template"
)

// RenderHome writes docs/index.mdx (site root) listing category cards from sidebar data.
func RenderHome(categories []SidebarCategory, out string) {
	t := template.Must(template.ParseFiles("templates/home.tpl"))
	f, err := os.Create(out)
	if err != nil {
		return
	}
	defer f.Close()
	_ = t.Execute(f, categories)
}
