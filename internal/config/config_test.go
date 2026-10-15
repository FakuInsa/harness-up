package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FakuInsa/harness-up/internal/config"
)

func TestConfig_Default(t *testing.T) {
	cfg := config.Default("MyProject")
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}

	if cfg.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", cfg.Version)
	}
	if cfg.Project.Name != "MyProject" {
		t.Errorf("expected project name MyProject, got %s", cfg.Project.Name)
	}
	if len(cfg.Targets) == 0 {
		t.Errorf("expected at least one default target, got none")
	}
	if cfg.Guardrails.MaxDiffLines != 300 {
		t.Errorf("expected max_diff_lines 300, got %d", cfg.Guardrails.MaxDiffLines)
	}
	if cfg.Memory.Provider != "sqlite-fts5" {
		t.Errorf("expected memory provider sqlite-fts5, got %s", cfg.Memory.Provider)
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got: %v", err)
	}
}

func TestConfig_RoundTrip(t *testing.T) {
	original := &config.Config{
		Version: "1.0",
		Project: config.ProjectConfig{
			Name: "Lowcord",
			Stack: config.StackConfig{
				Language:   "csharp",
				Framework:  "dotnet-9",
				TestRunner: "xunit",
			},
		},
		Targets: []string{"antigravity", "cursor"},
		Architecture: config.ArchitectureConfig{
			Pattern:          "vertical-slice",
			StrictBoundaries: true,
		},
		Workflows: []string{"sdd", "tdd-strict", "atomic-commits"},
		Conventions: config.ConventionsConfig{
			Principles: []string{"solid-strict", "immutability"},
			Naming:     "pascal-case",
		},
		Memory: config.MemoryConfig{
			Provider: "sqlite-fts5",
			Path:     ".harness/memory.db",
		},
		Guardrails: config.GuardrailsConfig{
			BlockDestructiveCommands: true,
			MaxDiffLines:             300,
			ProtectedPaths:           []string{".git/**", ".env*"},
		},
	}

	yamlData, err := original.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	loaded, err := config.LoadBytes(yamlData)
	if err != nil {
		t.Fatalf("LoadBytes failed: %v", err)
	}

	if loaded.Version != original.Version {
		t.Errorf("expected version %s, got %s", original.Version, loaded.Version)
	}
	if loaded.Project.Name != original.Project.Name {
		t.Errorf("expected project name %s, got %s", original.Project.Name, loaded.Project.Name)
	}
	if loaded.Project.Stack.Language != original.Project.Stack.Language {
		t.Errorf("expected stack language %s, got %s", original.Project.Stack.Language, loaded.Project.Stack.Language)
	}
	if loaded.Architecture.Pattern != original.Architecture.Pattern {
		t.Errorf("expected arch pattern %s, got %s", original.Architecture.Pattern, loaded.Architecture.Pattern)
	}
	if loaded.Guardrails.MaxDiffLines != original.Guardrails.MaxDiffLines {
		t.Errorf("expected max diff %d, got %d", original.Guardrails.MaxDiffLines, loaded.Guardrails.MaxDiffLines)
	}
	if len(loaded.Guardrails.ProtectedPaths) != 2 {
		t.Errorf("expected 2 protected paths, got %d", len(loaded.Guardrails.ProtectedPaths))
	}
}

func TestConfig_Validate_Errors(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(cfg *config.Config)
		errContains string
	}{
		{
			name: "empty version",
			modify: func(c *config.Config) {
				c.Version = ""
			},
			errContains: "version",
		},
		{
			name: "unsupported version",
			modify: func(c *config.Config) {
				c.Version = "2.0"
			},
			errContains: "unsupported version",
		},
		{
			name: "empty project name",
			modify: func(c *config.Config) {
				c.Project.Name = ""
			},
			errContains: "project.name",
		},
		{
			name: "empty targets",
			modify: func(c *config.Config) {
				c.Targets = []string{}
			},
			errContains: "targets",
		},
		{
			name: "unsupported target",
			modify: func(c *config.Config) {
				c.Targets = []string{"unsupported-agent"}
			},
			errContains: "target",
		},
		{
			name: "empty architecture pattern",
			modify: func(c *config.Config) {
				c.Architecture.Pattern = ""
			},
			errContains: "architecture.pattern",
		},
		{
			name: "unsupported memory provider",
			modify: func(c *config.Config) {
				c.Memory.Provider = "redis"
			},
			errContains: "memory.provider",
		},
		{
			name: "invalid max diff lines",
			modify: func(c *config.Config) {
				c.Guardrails.MaxDiffLines = 0
			},
			errContains: "max_diff_lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default("TestProject")
			tt.modify(cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatalf("expected validation error for %s, got nil", tt.name)
			}
			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errContains)) {
				t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
			}
		})
	}
}

func TestConfig_Save_And_Load_File(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "harnessup.yaml")

	original := config.Default("FileProject")
	original.Targets = []string{"antigravity", "cursor"}

	if err := config.Save(original, configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was physically written
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("config file is empty")
	}

	loaded, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Project.Name != "FileProject" {
		t.Errorf("expected FileProject, got %s", loaded.Project.Name)
	}
	if len(loaded.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(loaded.Targets))
	}
}

func TestConfig_Load_NonExistent(t *testing.T) {
	_, err := config.Load("non-existent-file-path-12345.yaml")
	if err == nil {
		t.Fatal("expected error loading non-existent file, got nil")
	}
}

func TestConfig_LoadBytes_InvalidYAML(t *testing.T) {
	_, err := config.LoadBytes([]byte("invalid: [yaml: unclosed"))
	if err == nil {
		t.Fatal("expected error parsing broken YAML, got nil")
	}
}

func TestConfig_Save_InvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "invalid.yaml")
	err := config.Save(&config.Config{}, path)
	if err == nil {
		t.Fatal("expected error saving invalid config, got nil")
	}
}

func TestConfig_Default_EmptyName(t *testing.T) {
	cfg := config.Default("")
	if cfg.Project.Name != "HarnessApp" {
		t.Errorf("expected fallback name HarnessApp, got %s", cfg.Project.Name)
	}
}

func TestConfig_Validate_MemoryProviders(t *testing.T) {
	providers := []string{"sqlite-fts5", "markdown-flat", "none"}
	for _, p := range providers {
		cfg := config.Default("MemTest")
		cfg.Memory.Provider = p
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected memory provider %q to be valid, got: %v", p, err)
		}
	}
}
