package config

import (
	"fmt"
	"strings"
)

var supportedTargets = map[string]bool{
	"antigravity":    true,
	"cursor":         true,
	"claude-code":    true,
	"opencode":       true,
	"windsurf":       true,
	"vscode-copilot": true,
	"pi":             true,
	"all":            true,
}

var supportedMemoryProviders = map[string]bool{
	"sqlite-fts5":   true,
	"markdown-flat": true,
	"none":          true,
}

// Validate verifica que la configuración cumpla con las invariantes del arnés.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Version) == "" {
		return fmt.Errorf("config validation failed: version is required")
	}
	if c.Version != "1.0" {
		return fmt.Errorf("config validation failed: unsupported version %q, expected '1.0'", c.Version)
	}

	if strings.TrimSpace(c.Project.Name) == "" {
		return fmt.Errorf("config validation failed: project.name cannot be empty")
	}

	if len(c.Targets) == 0 {
		return fmt.Errorf("config validation failed: targets must contain at least one agent target")
	}

	for _, target := range c.Targets {
		norm := strings.ToLower(strings.TrimSpace(target))
		if !supportedTargets[norm] {
			return fmt.Errorf("config validation failed: unsupported target %q", target)
		}
	}

	if strings.TrimSpace(c.Architecture.Pattern) == "" {
		return fmt.Errorf("config validation failed: architecture.pattern cannot be empty")
	}

	memNorm := strings.ToLower(strings.TrimSpace(c.Memory.Provider))
	if !supportedMemoryProviders[memNorm] {
		return fmt.Errorf("config validation failed: unsupported memory.provider %q, expected 'sqlite-fts5', 'markdown-flat' or 'none'", c.Memory.Provider)
	}

	if c.Guardrails.MaxDiffLines <= 0 {
		return fmt.Errorf("config validation failed: guardrails.max_diff_lines must be greater than 0, got %d", c.Guardrails.MaxDiffLines)
	}

	return nil
}
