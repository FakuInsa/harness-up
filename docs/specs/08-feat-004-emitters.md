# SPEC-008: Especificación Técnica — Emisores de Archivos Antigravity y Cursor (FEAT-004)

* **Proyecto:** HarnessUp
* **Feature ID:** `FEAT-004`
* **Estado:** Completado (Verificado con Evidencia)
* **Fecha:** 2026-10-15
* **Metodología:** Spec-Driven Development (SDD) & TDD Estricto
* **Rama Planificada:** `feat/feat-004-file-emitters`

---

## 1. Resumen y Objetivos

El objetivo de esta unidad de trabajo es implementar el subsistema de **Emisores de Archivos** (`internal/emitter`):
1. Definir la interfaz universal `Emitter` y el registro de emisores (`Registry`) para soportar múltiples agentes de IA de forma extensible.
2. Implementar el emisor específico para **Google Antigravity** (`AntigravityEmitter`), generando la estructura `.agent/rules/` y `.agent/skills/`.
3. Implementar el emisor específico para **Cursor** (`CursorEmitter`), generando archivos `.cursor/rules/*.mdc` con cabeceras frontmatter reconocidas por Cursor IDE.
4. Implementar el motor de fusión e idempotencia (`merge.go`), asegurando que la re-ejecución del arnés (`harnessup sync` o `init`) respete las personalizaciones del usuario delimitadas en zonas protegidas (`<!-- HARNESSUP:CUSTOM_USER_RULES -->`).

---

## 2. Contratos y Estructuras de Datos

### 2.1 Interfaz `Emitter` y Registro (`internal/emitter/emitter.go`)

```go
package emitter

import (
    "github.com/FakuInsa/harness-up/internal/config"
    "github.com/FakuInsa/harness-up/internal/engine"
)

// Emitter define el contrato para cualquier adaptador que genera artefactos de un agente de IA.
type Emitter interface {
    // ID retorna el identificador único del destino (ej. "antigravity", "cursor").
    ID() string
    // Emit genera físicamente los archivos del arnés en el directorio destino según la configuración.
    Emit(cfg *config.Config, r *engine.Renderer, targetDir string) error
}

// Registry gestiona la colección de emisores disponibles.
type Registry struct {
    emitters map[string]Emitter
}

// NewDefaultRegistry construye un registro con los emisores oficiales (Antigravity, Cursor).
func NewDefaultRegistry() *Registry

// Register asocia un emisor al registro.
func (reg *Registry) Register(e Emitter)

// Get obtiene un emisor por su ID.
func (reg *Registry) Get(id string) (Emitter, bool)

// EmitTargets ejecuta la emisión para todos los targets definidos en la configuración.
func (reg *Registry) EmitTargets(cfg *config.Config, r *engine.Renderer, targetDir string) error
```

---

### 2.2 Emisor Antigravity (`internal/emitter/antigravity.go`)

Genera la estructura de carpetas canónica para Google Antigravity:
```text
<targetDir>/
└── .agent/
    ├── rules/
    │   ├── architecture.md
    │   └── guardrails.md
    └── skills/
        └── sdd-workflow/
            └── SKILL.md
```

* Renderiza `rules/architecture.md.tmpl` $\rightarrow$ `.agent/rules/architecture.md`.
* Renderiza `rules/guardrails.md.tmpl` $\rightarrow$ `.agent/rules/guardrails.md`.
* Renderiza `skills/sdd-workflow/SKILL.md.tmpl` $\rightarrow$ `.agent/skills/sdd-workflow/SKILL.md`.

---

### 2.3 Emisor Cursor (`internal/emitter/cursor.go`)

Genera las reglas específicas para Cursor en formato `.mdc`:
```text
<targetDir>/
└── .cursor/
    └── rules/
        ├── architecture.mdc
        └── guardrails.mdc
```

Cada archivo `.mdc` incluye la cabecera frontmatter requerida por Cursor:
```markdown
---
description: {{.Description}}
globs: *
alwaysApply: true
---

{{.Content}}
```

---

### 2.4 Motor de Idempotencia y Fusión (`internal/emitter/merge.go`)

Para evitar la pérdida de reglas personalizadas escritas por el desarrollador cuando se re-ejecuta `harnessup init` o `harnessup sync`, todo archivo emitido utiliza el siguiente estándar de bloques:

```markdown
<!-- HARNESSUP:MANAGED_START -->
<!-- No edites este bloque manualmente. Es gestionado por HarnessUp. -->
... contenido renderizado automáticamente ...
<!-- HARNESSUP:MANAGED_END -->

<!-- HARNESSUP:CUSTOM_USER_RULES -->
... Cualquier regla o excepción personalizada agregada por el usuario aquí se preserva intacta ...
```

**Reglas de fusión:**
1. **Archivo nuevo:** Se genera con el bloque gestionado y el bloque de usuario vacío con comentarios guía.
2. **Archivo existente con bloques:** Se actualiza exclusivamente el contenido entre `MANAGED_START` y `MANAGED_END`, manteniendo 100% intacto el contenido existente en `CUSTOM_USER_RULES`.
3. **Archivo existente sin delimitadores (preexistente):** Se preserva el contenido previo intacto debajo del bloque gestionado para no provocar pérdida de datos.

---

## 3. Dependencias Externas

* **Cero dependencias externas:** Utilizar exclusivamente paquetes estándar de Go (`os`, `path/filepath`, `strings`, `fmt`, `bytes`).

---

## 4. Estrategia de TDD (Rojo $\rightarrow$ Verde $\rightarrow$ Refactor)

1. **Apertura de Rama:** `git checkout -b feat/feat-004-file-emitters` tras la aprobación del Gate.
2. **Paso RED:**
   * Crear `internal/emitter/emitter_test.go` cubriendo:
     * `TestAntigravityEmitter_Emit`: Verifica creación de `.agent/rules/*.md` y `.agent/skills/sdd-workflow/SKILL.md` en `t.TempDir()`.
     * `TestCursorEmitter_Emit`: Verifica creación de `.cursor/rules/*.mdc` con frontmatter válido.
     * `TestMerge_PreservesCustomUserRules`: Verifica que al re-emitir sobre un archivo existente, el bloque de usuario no se pierda ni se altere.
     * `TestMerge_PreexistingFileWithoutDelimiters`: Verifica que un archivo previo sin marcas conserve su contenido tras la emisión.
     * `TestRegistry_EmitTargets`: Verifica ejecución en cadena para múltiples targets (`targets: ["antigravity", "cursor"]`).
     * `TestRegistry_UnknownTarget`: Manejo de error claro ante un target no registrado.
   * Ejecutar `go test ./internal/emitter/...` y certificar falla inicial (RED).
3. **Paso GREEN:**
   * Implementar `internal/emitter/emitter.go`, `internal/emitter/antigravity.go`, `internal/emitter/cursor.go` y `internal/emitter/merge.go`.
   * Verificar que todos los tests pasen (GREEN).
4. **Paso REFACTOR & RECEIPT:**
   * Cobertura $\ge 90\%$.
   * `go vet ./...` y `gofmt` limpios.
   * Diff de código productivo $< 300$ líneas.

---

## 5. Criterios de Aceptación (DoD)

* [ ] Interfaz `Emitter` y estructura `Registry` implementadas en `internal/emitter`.
* [ ] `AntigravityEmitter` genera árbol `.agent/` válido y completo.
* [ ] `CursorEmitter` genera archivos `.cursor/rules/*.mdc` con frontmatter YAML válido.
* [ ] Fusión idempotente certificada: las modificaciones del usuario se preservan al re-emitir.
* [ ] 100% de tests unitarios pasando con cobertura $\ge 90\%$.
* [ ] Cero dependencias externas adicionales en `go.mod`.
* [ ] Diff productivo menor a 300 líneas.
