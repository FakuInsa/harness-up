package config

// Config representa el archivo maestro harnessup.yaml.
type Config struct {
	Version      string             `yaml:"version"`
	Project      ProjectConfig      `yaml:"project"`
	Targets      []string           `yaml:"targets"`
	Architecture ArchitectureConfig `yaml:"architecture"`
	Workflows    []string           `yaml:"workflows"`
	Conventions  ConventionsConfig  `yaml:"conventions"`
	Memory       MemoryConfig       `yaml:"memory"`
	Guardrails   GuardrailsConfig   `yaml:"guardrails"`
}

// ProjectConfig define metadatos y stack del proyecto.
type ProjectConfig struct {
	Name  string      `yaml:"name"`
	Stack StackConfig `yaml:"stack"`
}

// StackConfig define el lenguaje, framework y test runner.
type StackConfig struct {
	Language   string `yaml:"language"`
	Framework  string `yaml:"framework,omitempty"`
	TestRunner string `yaml:"test_runner,omitempty"`
}

// ArchitectureConfig define el patrón y fronteras de arquitectura.
type ArchitectureConfig struct {
	Pattern          string `yaml:"pattern"`
	StrictBoundaries bool   `yaml:"strict_boundaries"`
}

// ConventionsConfig define principios y estándares de nomenclatura.
type ConventionsConfig struct {
	Principles []string `yaml:"principles,omitempty"`
	Naming     string   `yaml:"naming,omitempty"`
}

// MemoryConfig define el proveedor y la ruta de persistencia técnica.
type MemoryConfig struct {
	Provider string `yaml:"provider"`
	Path     string `yaml:"path,omitempty"`
}

// GuardrailsConfig define reglas de seguridad y límites de diff.
type GuardrailsConfig struct {
	BlockDestructiveCommands bool     `yaml:"block_destructive_commands"`
	MaxDiffLines             int      `yaml:"max_diff_lines"`
	ProtectedPaths           []string `yaml:"protected_paths,omitempty"`
}
