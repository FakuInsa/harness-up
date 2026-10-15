package emitter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FakuInsa/harness-up/assets"
	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/emitter"
	"github.com/FakuInsa/harness-up/internal/engine"
)

func createTestConfig() *config.Config {
	cfg := config.Default("AcmeApp")
	cfg.Architecture.Pattern = "vertical-slice"
	cfg.Architecture.StrictBoundaries = true
	cfg.Conventions.Naming = "snake_case"
	cfg.Guardrails.MaxDiffLines = 300
	cfg.Guardrails.BlockDestructiveCommands = true
	cfg.Guardrails.ProtectedPaths = []string{".git/**", ".env*"}
	cfg.Targets = []string{"antigravity", "cursor"}
	return cfg
}

func TestAntigravityEmitter_Emit(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := createTestConfig()
	renderer := engine.NewRenderer(assets.FS)

	em := emitter.NewAntigravityEmitter()
	if em.ID() != "antigravity" {
		t.Fatalf("expected ID 'antigravity', got %q", em.ID())
	}

	if err := em.Emit(cfg, renderer, tmpDir); err != nil {
		t.Fatalf("unexpected emit error: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, ".agent", "rules", "architecture.md"),
		filepath.Join(tmpDir, ".agent", "rules", "guardrails.md"),
		filepath.Join(tmpDir, ".agent", "skills", "sdd-workflow", "SKILL.md"),
	}

	for _, path := range expectedFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read expected file %s: %v", path, err)
		}
		str := string(content)
		if !strings.Contains(str, "AcmeApp") {
			t.Errorf("file %s does not contain rendered project name 'AcmeApp'", path)
		}
		if !strings.Contains(str, emitter.ManagedStart) || !strings.Contains(str, emitter.ManagedEnd) {
			t.Errorf("file %s is missing managed delimiters", path)
		}
		if !strings.Contains(str, emitter.CustomUserRulesTag) {
			t.Errorf("file %s is missing custom user rules delimiter", path)
		}
	}
}

func TestCursorEmitter_Emit(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := createTestConfig()
	renderer := engine.NewRenderer(assets.FS)

	em := emitter.NewCursorEmitter()
	if em.ID() != "cursor" {
		t.Fatalf("expected ID 'cursor', got %q", em.ID())
	}

	if err := em.Emit(cfg, renderer, tmpDir); err != nil {
		t.Fatalf("unexpected emit error: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, ".cursor", "rules", "architecture.mdc"),
		filepath.Join(tmpDir, ".cursor", "rules", "guardrails.mdc"),
	}

	for _, path := range expectedFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read expected file %s: %v", path, err)
		}
		str := string(content)
		if !strings.HasPrefix(str, "---\n") {
			t.Errorf("file %s does not start with frontmatter delimiter", path)
		}
		if !strings.Contains(str, "globs: *") || !strings.Contains(str, "alwaysApply: true") {
			t.Errorf("file %s missing cursor frontmatter fields", path)
		}
		if !strings.Contains(str, emitter.ManagedStart) || !strings.Contains(str, emitter.ManagedEnd) {
			t.Errorf("file %s is missing managed delimiters", path)
		}
	}
}

func TestMerge_PreservesCustomUserRules(t *testing.T) {
	existingContent := strings.Join([]string{
		emitter.ManagedStart,
		"Old managed content",
		emitter.ManagedEnd,
		"",
		emitter.CustomUserRulesTag,
		"- Rule 1: Always check DB migrations",
		"- Rule 2: Never expose raw passwords",
	}, "\n")

	newManagedContent := "New managed content v2"

	merged := emitter.MergeContent(existingContent, newManagedContent)

	if !strings.Contains(merged, "New managed content v2") {
		t.Errorf("expected merged content to have new managed content, got: %s", merged)
	}
	if strings.Contains(merged, "Old managed content") {
		t.Errorf("merged content should not contain old managed content")
	}
	if !strings.Contains(merged, "- Rule 1: Always check DB migrations") ||
		!strings.Contains(merged, "- Rule 2: Never expose raw passwords") {
		t.Errorf("custom user rules were not preserved: %s", merged)
	}
}

func TestMerge_PreexistingFileWithoutDelimiters(t *testing.T) {
	existingContent := "# Preexisting manual file\nSome custom instructions."
	newManagedContent := "Generated managed block"

	merged := emitter.MergeContent(existingContent, newManagedContent)

	if !strings.Contains(merged, "Generated managed block") {
		t.Errorf("expected merged content to contain new managed content")
	}
	if !strings.Contains(merged, "# Preexisting manual file\nSome custom instructions.") {
		t.Errorf("expected legacy content to be preserved")
	}
	if !strings.Contains(merged, emitter.ManagedStart) || !strings.Contains(merged, emitter.CustomUserRulesTag) {
		t.Errorf("expected merged delimiters to be properly placed")
	}
}

func TestRegistry_EmitTargets(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := createTestConfig()
	cfg.Targets = []string{"antigravity", "cursor"}
	renderer := engine.NewRenderer(assets.FS)

	reg := emitter.NewDefaultRegistry()

	if err := reg.EmitTargets(cfg, renderer, tmpDir); err != nil {
		t.Fatalf("unexpected error emitting targets: %v", err)
	}

	antigravityArch := filepath.Join(tmpDir, ".agent", "rules", "architecture.md")
	cursorArch := filepath.Join(tmpDir, ".cursor", "rules", "architecture.mdc")

	if _, err := os.Stat(antigravityArch); err != nil {
		t.Errorf("expected antigravity architecture rule to exist: %v", err)
	}
	if _, err := os.Stat(cursorArch); err != nil {
		t.Errorf("expected cursor architecture rule to exist: %v", err)
	}
}

func TestRegistry_UnknownTarget(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := createTestConfig()
	cfg.Targets = []string{"unknown-agent"}
	renderer := engine.NewRenderer(assets.FS)

	reg := emitter.NewDefaultRegistry()

	err := reg.EmitTargets(cfg, renderer, tmpDir)
	if err == nil {
		t.Fatalf("expected error for unknown target, got nil")
	}
	if !strings.Contains(err.Error(), "unknown target") {
		t.Errorf("expected error message to mention unknown target, got %v", err)
	}
}

func TestMerge_ExistingManagedWithoutUserTag(t *testing.T) {
	existing := strings.Join([]string{
		emitter.ManagedStart,
		"Managed body",
		emitter.ManagedEnd,
	}, "\n")

	merged := emitter.MergeContent(existing, "Updated body")
	if !strings.Contains(merged, "Updated body") {
		t.Errorf("expected updated body")
	}
	if !strings.Contains(merged, emitter.CustomUserRulesTag) {
		t.Errorf("expected custom user rules tag to be appended when missing")
	}
}

func TestRegistry_Get(t *testing.T) {
	reg := emitter.NewDefaultRegistry()
	em, ok := reg.Get("antigravity")
	if !ok || em == nil {
		t.Fatalf("expected to find antigravity emitter")
	}
	if em.ID() != "antigravity" {
		t.Errorf("expected ID 'antigravity', got %q", em.ID())
	}

	_, ok = reg.Get("nonexistent")
	if ok {
		t.Errorf("expected nonexistent to return false")
	}
}

type errorEmitter struct{}

func (e *errorEmitter) ID() string { return "faulty" }
func (e *errorEmitter) Emit(_ *config.Config, _ *engine.Renderer, _ string) error {
	return os.ErrPermission
}

func TestRegistry_EmitTargets_EmitterError(t *testing.T) {
	reg := emitter.NewDefaultRegistry()
	reg.Register(&errorEmitter{})
	cfg := createTestConfig()
	cfg.Targets = []string{"faulty"}
	renderer := engine.NewRenderer(assets.FS)

	err := reg.EmitTargets(cfg, renderer, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "failed to emit target") {
		t.Fatalf("expected failed to emit target error, got: %v", err)
	}
}

func TestEmitters_RenderAssetError(t *testing.T) {
	// Renderer with an empty FS will fail on missing template
	emptyRenderer := engine.NewRenderer(os.DirFS(t.TempDir()))
	cfg := createTestConfig()
	tmpDir := t.TempDir()

	antigravity := emitter.NewAntigravityEmitter()
	if err := antigravity.Emit(cfg, emptyRenderer, tmpDir); err == nil {
		t.Errorf("expected error when asset is missing in antigravity emitter")
	}

	cursor := emitter.NewCursorEmitter()
	if err := cursor.Emit(cfg, emptyRenderer, tmpDir); err == nil {
		t.Errorf("expected error when asset is missing in cursor emitter")
	}
}


