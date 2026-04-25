// @ts-check
const sidebars = {
  tutorialSidebar: [
    {
      type: 'doc',
      id: 'index',
      label: 'Home',
    },
{{- range . }}
    {
      type: 'category',
      label: '{{.Name}}',
      link: {
        type: 'doc',
        id: '{{.Slug}}/index',
      },
      items: [
{{- if .Concepts }}
        {
          type: 'category',
          label: '📘 Concepts',
          items: [
{{- range .Concepts }}
            '{{.}}',
{{- end }}
          ],
        },
{{- end }}
{{- if .Problems }}
        {
          type: 'category',
          label: '🧩 Problems',
          items: [
{{- range .Problems }}
            '{{.}}',
{{- end }}
          ],
        },
{{- end }}
      ],
    },
{{- end }}
  ],
};

export default sidebars;
