package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Default construye una configuración canónica con convenciones recomendadas.
func Default(projectName string) *Config {
	if projectName == "" {
		projectName = "HarnessApp"
	}
	return &Config{
		Version: "1.0",
		Project: ProjectConfig{
			Name: projectName,
			Stack: StackConfig{
				Language: "unknown",
			},
		},
		Targets: []string{"antigravity"},
		Architecture: ArchitectureConfig{
			Pattern:          "vertical-slice",
			StrictBoundaries: true,
		},
		Workflows: []string{"sdd", "tdd-strict", "atomic-commits"},
		Conventions: ConventionsConfig{
			Principles: []string{"solid-strict", "immutability"},
			Naming:     "standard",
		},
		Memory: MemoryConfig{
			Provider: "sqlite-fts5",
			Path:     ".harness/memory.db",
		},
		Guardrails: GuardrailsConfig{
			BlockDestructiveCommands: true,
			MaxDiffLines:             300,
			ProtectedPaths: []string{
				".git/**",
				".env*",
			},
		},
	}
}

// ToYAML serializa la configuración a YAML determinista.
func (c *Config) ToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

// LoadBytes deserializa un slice de bytes YAML y valida el arnés.
func LoadBytes(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Load lee el archivo YAML en la ruta indicada, lo parsea y lo valida.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}
	return LoadBytes(data)
}

// Save valida la configuración y la persiste en disco con permisos 0644.
func Save(cfg *Config, path string) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("cannot save invalid config: %w", err)
	}

	data, err := cfg.ToYAML()
	if err != nil {
		return fmt.Errorf("failed to serialize config to yaml: %w", err)
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dir, err)
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %q: %w", path, err)
	}

	return nil
}
