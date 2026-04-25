package renderer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WriteGithubRoutes writes a JSON map: Docusaurus doc id -> GitHub URL (blob or tree).
func WriteGithubRoutes(outPath string, routes map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, append(b, '\n'), 0o644)
}
