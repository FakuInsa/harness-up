# SPEC-006: Especificación Técnica — Modelo de Configuración & Serialización (FEAT-002)

* **Proyecto:** HarnessUp
* **Feature ID:** `FEAT-002`
* **Estado:** Aprobado por el Tech Lead
* **Fecha:** 2026-10-15
* **Metodología:** Spec-Driven Development (SDD) & TDD Estricto
* **Rama Planificada:** `feat/feat-002-config-model`

---

## 1. Resumen y Objetivos

El objetivo de esta unidad de trabajo es implementar el paquete `internal/config`, responsable de:
1. Definir los structs canónicos de Go que representan el estado completo de un arnés según lo estipulado en [SPEC-004](./04-schemas-and-emitters.md).
2. Serializar y deserializar el archivo declarativo `harnessup.yaml` con idempotencia estricta.
3. Proveer validación de esquemas con mensajes de error descriptivos ante entradas inválidas o incompletas.
4. Proveer una configuración base por defecto (`Default()`) para inicializaciones limpias.

---

## 2. Contratos y Estructuras de Datos (`internal/config`)

### 2.1 Structs de Configuración

```go
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

type ProjectConfig struct {
    Name  string      `yaml:"name"`
    Stack StackConfig `yaml:"stack"`
}

type StackConfig struct {
    Language   string `yaml:"language"`
    Framework  string `yaml:"framework,omitempty"`
    TestRunner string `yaml:"test_runner,omitempty"`
}

type ArchitectureConfig struct {
    Pattern          string `yaml:"pattern"`
    StrictBoundaries bool   `yaml:"strict_boundaries"`
}

type ConventionsConfig struct {
    Principles []string `yaml:"principles,omitempty"`
    Naming     string   `yaml:"naming,omitempty"`
}

type MemoryConfig struct {
    Provider string `yaml:"provider"`
    Path     string `yaml:"path,omitempty"`
}

type GuardrailsConfig struct {
    BlockDestructiveCommands bool     `yaml:"block_destructive_commands"`
    MaxDiffLines             int      `yaml:"max_diff_lines"`
    ProtectedPaths           []string `yaml:"protected_paths,omitempty"`
}
```

### 2.2 Funciones y Métodos del Paquete

```go
// Default genera una configuración lista para usar con valores razonables de la industria.
func Default(projectName string) *Config

// Load lee y valida un archivo harnessup.yaml desde la ruta indicada.
func Load(path string) (*Config, error)

// LoadBytes parsea y valida un slice de bytes YAML.
func LoadBytes(data []byte) (*Config, error)

// Save serializa la configuración y la escribe en el archivo destino con permisos 0644.
func Save(cfg *Config, path string) error

// ToYAML serializa la configuración a formato YAML en memoria.
func (c *Config) ToYAML() ([]byte, error)

// Validate verifica que la configuración no tenga campos vacíos o valores no soportados.
func (c *Config) Validate() error
```

---

## 3. Reglas de Validación Declarativa

El método `Validate()` aplicará las siguientes invariantes:
1. `Version`: Obligatoria (debe ser `"1.0"`).
2. `Project.Name`: No puede estar vacío.
3. `Targets`: Debe contener al menos un target soportado (`antigravity`, `cursor`, `claude-code`, `opencode`, `windsurf`, `vscode-copilot`, `pi`, o `all`).
4. `Architecture.Pattern`: No puede estar vacío (ej. `vertical-slice`, `clean`, `hexagonal`, `modular-monolith`).
5. `Memory.Provider`: Debe ser uno de: `sqlite-fts5`, `markdown-flat`, o `none`.
6. `Guardrails.MaxDiffLines`: Debe ser estrictamente mayor a 0 (por defecto 300).

---

## 4. Dependencias Externas

* Conforme al **Protocolo de Interrupción (Punto 3)** del [Way of Working](../WAY_OF_WORKING.md#5-protocolo-de-interrupción-cuándo-debe-frenar-la-ia), se propone la incorporación de:
  * `gopkg.in/yaml.v3` (estándar de la industria en Go para serialización YAML con preservación de formato).

---

## 5. Estrategia de TDD (Rojo $\rightarrow$ Verde $\rightarrow$ Refactor)

1. **Apertura de Rama:** `git checkout -b feat/feat-002-config-model` tras la aprobación del Gate.
2. **Paso RED:**
   * Crear `internal/config/config_test.go`.
   * Pruebas a implementar:
     * `TestConfig_Default`: Verifica valores sensatos por defecto.
     * `TestConfig_RoundTrip`: Serializa un `Config` a YAML, lo vuelve a cargar (`LoadBytes`), y certifica igualdad exacta sin pérdida de campos.
     * `TestConfig_Validate_Success`: Valida un YAML completo como el presentado en SPEC-004.
     * `TestConfig_Validate_Errors`: Comprueba que se rechacen configuraciones con:
       * Nombre de proyecto vacío.
       * Targets vacíos o desconocidos.
       * Versión no soportada.
       * Proveedor de memoria inválido.
       * `max_diff_lines` $\le 0$.
     * `TestConfig_Save_And_Load_File`: Comprueba escritura y lectura física en disco temporal (`t.TempDir()`).
   * Ejecutar y verificar falla (RED).
3. **Paso GREEN:**
   * Agregar dependencia `gopkg.in/yaml.v3`.
   * Implementar `internal/config/config.go` y `internal/config/validate.go`.
   * Pasar los tests al 100% (GREEN).
4. **Paso REFACTOR & RECEIPT:**
   * Ejecutar `go test -v -cover ./internal/config/...`, `go vet ./...` y `gofmt`.
   * Cobertura objetivo $\ge 90\%$.
   * Diff de código productivo $< 300$ líneas.

---

## 6. Criterios de Aceptación (DoD)

* [ ] Dependencia `gopkg.in/yaml.v3` registrada limpiamente en `go.mod`.
* [ ] Paquete `internal/config` implementado según los contratos.
* [ ] 100% tests pasando con cobertura $\ge 90\%$.
* [ ] Cero advertencias en `go vet ./...`.
* [ ] Diff productivo menor a 300 líneas.
