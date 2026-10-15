package engine

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"
)

// Renderer encapsula el motor de renderizado basado en text/template y un fs.FS.
type Renderer struct {
	fs      fs.FS
	funcMap template.FuncMap
}

// NewRenderer construye un Renderer asociado a un sistema de archivos (ej. assets.FS o fstest.MapFS).
func NewRenderer(fileSystem fs.FS) *Renderer {
	return &Renderer{
		fs:      fileSystem,
		funcMap: defaultFuncMap(),
	}
}

// RenderString compila y ejecuta una plantilla en texto plano pasando un modelo de datos.
func (r *Renderer) RenderString(tmplContent string, data any) (string, error) {
	tmpl, err := template.New("string").Funcs(r.funcMap).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderAsset lee una plantilla desde el fs embebido, la compila y devuelve el texto procesado.
func (r *Renderer) RenderAsset(assetPath string, data any) (string, error) {
	content, err := fs.ReadFile(r.fs, assetPath)
	if err != nil {
		return "", fmt.Errorf("failed to read asset %q: %w", assetPath, err)
	}

	return r.RenderString(string(content), data)
}

// ListAssets lista las rutas de archivos disponibles dentro de un prefijo en el FS.
func (r *Renderer) ListAssets(prefix string) ([]string, error) {
	var files []string

	err := fs.WalkDir(r.fs, prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list assets in %q: %w", prefix, err)
	}

	return files, nil
}
