package emitter

import (
	"fmt"
	"path/filepath"

	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/engine"
)

// AntigravityEmitter genera artefactos de reglas y skills para Google Antigravity.
type AntigravityEmitter struct{}

// NewAntigravityEmitter construye una nueva instancia de AntigravityEmitter.
func NewAntigravityEmitter() *AntigravityEmitter {
	return &AntigravityEmitter{}
}

// ID retorna el identificador oficial del target.
func (e *AntigravityEmitter) ID() string {
	return "antigravity"
}

// Emit genera los archivos `.agent/rules/` y `.agent/skills/` aplicando fusión idempotente.
func (e *AntigravityEmitter) Emit(cfg *config.Config, r *engine.Renderer, targetDir string) error {
	archContent, err := r.RenderAsset("rules/architecture.md.tmpl", cfg)
	if err != nil {
		return fmt.Errorf("failed to render architecture rule: %w", err)
	}
	archPath := filepath.Join(targetDir, ".agent", "rules", "architecture.md")
	if err := WriteMergedFile(archPath, archContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", archPath, err)
	}

	guardContent, err := r.RenderAsset("rules/guardrails.md.tmpl", cfg)
	if err != nil {
		return fmt.Errorf("failed to render guardrails rule: %w", err)
	}
	guardPath := filepath.Join(targetDir, ".agent", "rules", "guardrails.md")
	if err := WriteMergedFile(guardPath, guardContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", guardPath, err)
	}

	skillContent, err := r.RenderAsset("skills/sdd-workflow/SKILL.md.tmpl", cfg)
	if err != nil {
		return fmt.Errorf("failed to render sdd-workflow skill: %w", err)
	}
	skillPath := filepath.Join(targetDir, ".agent", "skills", "sdd-workflow", "SKILL.md")
	if err := WriteMergedFile(skillPath, skillContent, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", skillPath, err)
	}

	return nil
}
