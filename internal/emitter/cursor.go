package emitter

import (
	"fmt"
	"path/filepath"

	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/engine"
)

// CursorEmitter genera reglas de contexto para Cursor IDE en formato .mdc con frontmatter YAML.
type CursorEmitter struct{}

// NewCursorEmitter construye una nueva instancia de CursorEmitter.
func NewCursorEmitter() *CursorEmitter {
	return &CursorEmitter{}
}

// ID retorna el identificador oficial del target.
func (e *CursorEmitter) ID() string {
	return "cursor"
}

// Emit genera los archivos `.cursor/rules/*.mdc` aplicando frontmatter y fusión idempotente.
func (e *CursorEmitter) Emit(cfg *config.Config, r *engine.Renderer, targetDir string) error {
	archContent, err := r.RenderAsset("rules/architecture.md.tmpl", cfg)
	if err != nil {
		return fmt.Errorf("failed to render architecture rule for cursor: %w", err)
	}
	archFrontmatter := fmt.Sprintf("---\ndescription: Reglas de arquitectura y convenciones para %s\nglobs: *\nalwaysApply: true\n---", cfg.Project.Name)
	archPath := filepath.Join(targetDir, ".cursor", "rules", "architecture.mdc")
	if err := WriteMergedFileWithHeader(archPath, archFrontmatter, archContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", archPath, err)
	}

	guardContent, err := r.RenderAsset("rules/guardrails.md.tmpl", cfg)
	if err != nil {
		return fmt.Errorf("failed to render guardrails rule for cursor: %w", err)
	}
	guardFrontmatter := fmt.Sprintf("---\ndescription: Guardarraíles de seguridad y calidad para %s\nglobs: *\nalwaysApply: true\n---", cfg.Project.Name)
	guardPath := filepath.Join(targetDir, ".cursor", "rules", "guardrails.mdc")
	if err := WriteMergedFileWithHeader(guardPath, guardFrontmatter, guardContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", guardPath, err)
	}

	return nil
}
