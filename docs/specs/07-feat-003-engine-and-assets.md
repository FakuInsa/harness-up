# SPEC-007: Especificación Técnica — Motor de Plantillas y Assets Embebidos (FEAT-003)

* **Proyecto:** HarnessUp
* **Feature ID:** `FEAT-003`
* **Estado:** [✓] Completado e Integrado (DONE)
* **Fecha:** 2026-10-15
* **Metodología:** Spec-Driven Development (SDD) & TDD Estricto
* **Rama Planificada:** `feat/feat-003-engine-and-assets`

---

## 1. Resumen y Objetivos

El objetivo de esta unidad de trabajo es implementar el subsistema de activos embebidos y compilación de plantillas:
1. Crear el paquete raíz `assets` con la directiva `//go:embed` para compilar dentro del binario Go los fragmentos de reglas, plantillas maestras y skills.
2. Implementar el paquete `internal/engine` con un compilador de plantillas (`text/template`) enriquecido con funciones auxiliares de transformación de texto (`toLower`, `toUpper`, `camelCase`, `pascalCase`, `kebabCase`, `snakeCase`).
3. Empaquetar las primeras plantillas de reglas (`architecture.md.tmpl`, `guardrails.md.tmpl`) y skills (`sdd-workflow/SKILL.md.tmpl`) listas para ser compiladas y emitidas en features posteriores.

---

## 2. Contratos y Estructuras de Datos

### 2.1 Paquete `assets` (`assets/assets.go`)

```go
package assets

import "embed"

// FS contiene todos los archivos estáticos de reglas, skills y templates embebidos en el binario.
//go:embed rules/* skills/* templates/*
var FS embed.FS
```

Estructura del directorio `assets/`:
```text
assets/
├── assets.go
├── rules/
│   ├── architecture.md.tmpl
│   └── guardrails.md.tmpl
├── skills/
│   └── sdd-workflow/
│       └── SKILL.md.tmpl
└── templates/
    └── .keep
```

### 2.2 Paquete `internal/engine` (`internal/engine/renderer.go`)

```go
package engine

import (
    "html/template"
    "io/fs"
)

// Renderer encapsula el motor de renderizado basado en text/template y un fs.FS.
type Renderer struct {
    fs fs.FS
}

// NewRenderer construye un Renderer asociado a un sistema de archivos (ej. assets.FS o fstest.MapFS).
func NewRenderer(fileSystem fs.FS) *Renderer

// RenderString compila y ejecuta una plantilla en texto plano pasando un modelo de datos.
func (r *Renderer) RenderString(tmplContent string, data any) (string, error)

// RenderAsset lee una plantilla desde el fs embebido, la compila y devuelve el texto procesado.
func (r *Renderer) RenderAsset(assetPath string, data any) (string, error)

// ListAssets lista las rutas de archivos disponibles dentro de un prefijo en el FS.
func (r *Renderer) ListAssets(prefix string) ([]string, error)
```

### 2.3 Funciones Auxiliares Disponibles en Plantillas (`FuncMap`)

* `toLower`: pasa cadena a minúsculas.
* `toUpper`: pasa cadena a mayúsculas.
* `camelCase`: convierte `hello-world` / `hello_world` $\rightarrow$ `helloWorld`.
* `pascalCase`: convierte `hello-world` / `hello_world` $\rightarrow$ `HelloWorld`.
* `snakeCase`: convierte `helloWorld` / `hello-world` $\rightarrow$ `hello_world`.
* `kebabCase`: convierte `helloWorld` / `hello_world` $\rightarrow$ `hello-world`.
* `contains`: evalúa si un slice de strings o una cadena contiene un elemento.
* `join`: une un slice de strings con un separador.

---

## 3. Plantillas Iniciales del Catálogo

### 3.1 `assets/rules/architecture.md.tmpl`
Genera las directivas de arquitectura según `{{.Architecture.Pattern}}`:
* Si es `vertical-slice`: prohíbe capas genéricas `services/`, `repositories/` y exige agrupación por feature.
* Si es `clean`: impone las 4 capas concéntricas (Domain, Application, Infrastructure, API).
* Si es `hexagonal`: impone puertos y adaptadores.

### 3.2 `assets/rules/guardrails.md.tmpl`
Genera las reglas de seguridad:
* Bloqueo de comandos destructivos (`rm -rf`, `git reset --hard`, etc.) si `{{.Guardrails.BlockDestructiveCommands}}` es verdadero.
* Límite de diff por unidad de trabajo (`{{.Guardrails.MaxDiffLines}}`).
* Rutas protegidas (`{{range .Guardrails.ProtectedPaths}}`).

---

## 4. Dependencias Externas

* **Cero dependencias externas:** Se utilizará exclusivamente la biblioteca estándar de Go (`io/fs`, `text/template`, `strings`, `bytes`, `unicode`).

---

## 5. Estrategia de TDD (Rojo $\rightarrow$ Verde $\rightarrow$ Refactor)

1. **Apertura de Rama:** `git checkout -b feat/feat-003-engine-and-assets` tras la aprobación del Gate.
2. **Paso RED:**
   * Crear `internal/engine/renderer_test.go` con tests:
     * `TestRenderer_RenderString_FuncMap`: Pruebas de `camelCase`, `pascalCase`, `snakeCase`, etc.
     * `TestRenderer_RenderAsset_Embedded`: Renderizado de `rules/architecture.md.tmpl` y `rules/guardrails.md.tmpl` pasando un `config.Config`.
     * `TestRenderer_Conditionals`: Evaluación de bloques condicionales (`{{if eq .Architecture.Pattern "vertical-slice"}}`).
     * `TestRenderer_MissingAsset`: Error claro al solicitar una plantilla inexistente.
     * `TestRenderer_InvalidSyntax`: Error ante sintaxis rota de Go template.
     * `TestRenderer_ListAssets`: Comprueba listado de activos por prefijo.
   * Ejecutar `go test ./internal/engine/...` y certificar falla inicial (RED).
3. **Paso GREEN:**
   * Crear `assets/assets.go` y las plantillas iniciales.
   * Implementar `internal/engine/renderer.go` y `internal/engine/funcs.go`.
   * Verificar que todos los tests pasen (GREEN).
4. **Paso REFACTOR & RECEIPT:**
   * Cobertura $\ge 90\%$.
   * `go vet ./...` y `gofmt` limpios.
   * Diff de código productivo $< 300$ líneas.

---

## 6. Criterios de Aceptación (DoD)

* [ ] Paquete `assets` compila correctamente con `//go:embed`.
* [ ] Paquete `internal/engine` implementado y desacoplado vía `fs.FS`.
* [ ] 100% de tests unitarios pasando con cobertura $\ge 90\%$.
* [ ] Cero dependencias externas adicionales en `go.mod`.
* [ ] Diff productivo menor a 300 líneas.
