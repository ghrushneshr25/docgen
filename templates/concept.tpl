---
title: {{.Title}}
---

# {{.Title}}

<br />

<table>
  <tbody>
    <tr><td>Difficulty</td><td>{{index .Meta "difficulty"}}</td></tr>
    <tr><td>Tags</td><td>{{index .Meta "tags"}}</td></tr>
  </tbody>
</table>

---

{{if .ConceptDescription}}
## Description

{{.ConceptDescription}}
{{end}}

## Structure

{{if .ConceptStructureIntro}}
{{.ConceptStructureIntro}}

{{end}}
{{range .Structs}}
### {{.Name}}

{{if .Doc}}
{{.Doc}}

{{end}}
```go
{{.Code}}
```

{{end}}

{{if .OperationSubsections}}
## Operations
{{range .OperationSubsections}}
### {{.Title}}
{{.Prose}}

```go
{{.Code}}
```

{{end}}
{{else if .ConceptOperationsMD}}
## Operations

{{.ConceptOperationsMD}}
{{else if .OperationsCode}}
## Operations

{{end}}
{{if and (eq (len .OperationSubsections) 0) .OperationsCode}}
```go
{{.OperationsCode}}
```

{{end}}
