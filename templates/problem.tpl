---
title: {{.Title}}
sidebar_position: {{.SidebarPosition}}
---

# {{.Title}}

<br />

<table>
  <tbody>
    <tr><td>Difficulty</td><td>{{index .Meta "difficulty"}}</td></tr>
    <tr><td>Tags</td><td>{{index .Meta "tags"}}</td></tr>
    <tr><td>Time complexity</td><td>{{index .Meta "time"}}</td></tr>
    <tr><td>Space complexity</td><td>{{index .Meta "space"}}</td></tr>
  </tbody>
</table>

---

{{.Sections}}

---

## 🧩 Solution (Go)

```go
{{.Code}}
```

{{if .HasTests}}

---

## 🧪 Tests

{{range .Subtests}}
#### {{.Name}}

{{if .Desc}}
{{.Desc}}

{{end}}
```go
{{.Code}}
```

{{end}}
{{end}}
