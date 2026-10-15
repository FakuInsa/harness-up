package detector

import "io/fs"

// Language identifica el lenguaje de programación detectado.
type Language string

const (
	LangCSharp     Language = "csharp"
	LangTypeScript Language = "typescript"
	LangPython     Language = "python"
	LangGo         Language = "go"
	LangUnknown    Language = "unknown"
)

// TestRunner identifica el framework de pruebas detectado.
type TestRunner string

const (
	RunnerXUnit   TestRunner = "xunit"
	RunnerNUnit   TestRunner = "nunit"
	RunnerMSTest  TestRunner = "mstest"
	RunnerVitest  TestRunner = "vitest"
	RunnerJest    TestRunner = "jest"
	RunnerPytest  TestRunner = "pytest"
	RunnerGoTest  TestRunner = "gotest"
	RunnerUnknown TestRunner = "unknown"
)

// StackInfo contiene el resultado del análisis del repositorio.
type StackInfo struct {
	Language   Language   `json:"language"`
	Framework  string     `json:"framework,omitempty"`
	TestRunner TestRunner `json:"test_runner"`
	Confidence float32    `json:"confidence"`
	Evidence   []string   `json:"evidence"`
}

// Detector define el contrato para el análisis de stack.
type Detector interface {
	DetectFS(sysFS fs.FS) (*StackInfo, error)
	DetectDir(dirPath string) (*StackInfo, error)
}
