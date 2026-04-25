---
title: {{.Title}}
slug: /{{.Slug}}
pagination_prev: null
pagination_next: null
---

# 📂 {{.Title}}

{{if .ConceptItems}}
## Concepts

Documentation pages for **concepts** in this category (data structures, APIs, theory).

| Concept | Difficulty | Tags |
| --- | --- | --- |
{{- range .ConceptItems }}
| [{{.Title}}]({{.Link}}) | {{.Difficulty}} | {{.Tags}} |
{{- end }}

{{end}}
{{- if and .ConceptItems .ProblemItems }}

---

{{ end }}
{{- if .ProblemItems }}
## Problems

**Practice problems** in this category (with solutions and tests).

| Problem | Difficulty | Tags |
| --- | --- | --- |
{{- range .ProblemItems }}
| [{{.Title}}]({{.Link}}) | {{.Difficulty}} | {{.Tags}} |
{{- end }}

{{end}}
