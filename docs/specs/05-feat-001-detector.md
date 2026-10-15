# SPEC-005: Especificación Técnica — Scaffolding Base & Detector de Stack (FEAT-001)

* **Proyecto:** HarnessUp
* **Feature ID:** `FEAT-001`
* **Estado:** Aprobado por el Tech Lead
* **Fecha:** 2026-10-15
* **Metodología:** Spec-Driven Development (SDD) & TDD Estricto

---

## 1. Resumen y Objetivos

El objetivo de esta unidad de trabajo es:
1. Inicializar el módulo Go base del proyecto (`github.com/FakuInsa/harness-up` o nombre de módulo acordado en `go.mod`).
2. Implementar el paquete `internal/detector`, encargado de analizar de manera determinista el árbol de archivos de un proyecto para inferir:
   * **Lenguaje de programación principal** (`csharp`, `typescript`, `python`, `go`, o `unknown`).
   * **Framework o runtime relevante** (ej. `dotnet-9`, `nodejs`, `fastapi`, etc.).
   * **Test runner existente** (ej. `xunit`, `nunit`, `vitest`, `jest`, `pytest`, `gotest`, o `unknown`).
   * **Archivos marcadores identificados** (evidencia documental del hallazgo).

---

## 2. Contratos y Estructuras de Datos (`internal/detector`)

### 2.1 Tipos y Enums

```go
package detector

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
    Language    Language   `json:"language"`
    Framework   string     `json:"framework,omitempty"`
    TestRunner  TestRunner `json:"test_runner"`
    Confidence  float32    `json:"confidence"` // 0.0 a 1.0
    Evidence    []string   `json:"evidence"`    // Rutas relativas que fundamentaron la decisión
}
```

### 2.2 Interfaz del Detector

Para facilitar el testing unitario y desacoplar de llamadas al disco físico, la interfaz soportará tanto un sistema de archivos abstracto (`fs.FS`) como una función de conveniencia sobre un directorio del SO:

```go
// Detector define el contrato para el análisis de stack.
type Detector interface {
    DetectFS(sysFS fs.FS) (*StackInfo, error)
    DetectDir(dirPath string) (*StackInfo, error)
}
```

---

## 3. Matriz de Detección Heurística

El detector escanea el directorio raíz (y hasta 2 niveles de profundidad para proyectos monorepo o subcarpetas comunes como `src/` o `apps/`) buscando los siguientes marcadores:

| Lenguaje | Marcadores de Archivo | Heurística de Test Runner | Heurística de Framework |
| :--- | :--- | :--- | :--- |
| **C#** | `*.sln`, `*.csproj` | Inspección de contenido de `.csproj` en busca de `PackageReference`: `xunit`, `NUnit`, `Microsoft.NET.Test.Sdk` | Inspección de `TargetFramework` (ej. `net9.0`, `net8.0`) |
| **TypeScript / JS** | `package.json`, `tsconfig.json` | Inspección de `package.json` (`devDependencies`, `dependencies`, o `scripts.test`): `vitest`, `jest`, `mocha` | Presencia de `@nestjs`, `next`, `react`, `express` |
| **Python** | `pyproject.toml`, `requirements.txt`, `Pipfile` | Inspección de `pyproject.toml` o `requirements.txt`: `pytest` | Presencia de `fastapi`, `django`, `flask` |
| **Go** | `go.mod` | Presencia de `_test.go` o dependencias como `testify` $\rightarrow$ `gotest` | Versión de `go` en `go.mod` |
| **Desconocido** | Ningún marcador reconocido | `RunnerUnknown` | Devuelve `LangUnknown` con `Confidence: 0.0` para que el Wizard TUI consulte al usuario. |

---

## 4. Estrategia de TDD (Rojo $\rightarrow$ Verde $\rightarrow$ Refactor)

Siguiendo el principio de **Test-First**:

1. **Paso RED:**
   * Crear `internal/detector/detector_test.go`.
   * Implementar casos de prueba usando `testing/fstest.MapFS` para simular en memoria estructuras de repositorios:
     * `TestDetect_CSharp_XUnit`: repositorio con `.sln` y archivo `.csproj` conteniendo `PackageReference Include="xunit"`.
     * `TestDetect_TypeScript_Vitest`: repositorio con `package.json` conteniendo `"vitest": "^2.0.0"`.
     * `TestDetect_Python_Pytest`: repositorio con `pyproject.toml` conteniendo `[tool.pytest.ini_options]`.
     * `TestDetect_Go_Standard`: repositorio con `go.mod`.
     * `TestDetect_EmptyOrUnknown`: directorio vacío o sin marcadores.
   * Ejecutar `go test ./internal/detector/...` y constatar la falla inicial (RED).

2. **Paso GREEN:**
   * Implementar `internal/detector/detector.go` con los analizadores mínimos necesarios para satisfacer los tests.
   * Ejecutar y verificar que los tests pasan (GREEN).

3. **Paso REFACTOR / EVIDENCE:**
   * Limpiar la lógica de parsing de archivos sin romper los tests.
   * Ejecutar `go test -v ./internal/detector/...` y documentar la evidencia.

---

## 5. Criterios de Aceptación (DoD de FEAT-001)

* [ ] Módulo Go inicializado con Go 1.22+.
* [ ] Paquete `internal/detector` implementado con cero dependencias externas pesadas (usar solo stdlib de Go: `io/fs`, `path/filepath`, `encoding/json`, `strings`, `bufio`).
* [ ] 100% de tests unitarios pasando en `internal/detector/detector_test.go`.
* [ ] Cobertura de tests $\ge$ 90% en el paquete `detector`.
* [ ] `go vet ./...` y `gofmt` ejecutados sin advertencias.
* [ ] El diff total es menor a 300 líneas.
