# SPEC-004: Esquemas de Datos, Motor de Plantillas y Emisión

* **Proyecto:** HarnessUp
* **Estado:** Propuesto (Pendiente de Aprobación)
* **Fecha:** 2026-10-14
* **Metodología:** Spec-Driven Development (SDD)

---

## 1. Esquema de Configuración del Proyecto (`harnessup.yaml`)

El resultado del wizard o la ejecución manual se guarda en la raíz del proyecto en un archivo declarativo `harnessup.yaml`:

```yaml
version: "1.0"
project:
  name: "MyProject"
  stack:
    language: "csharp"
    framework: "dotnet-9"
    test_runner: "xunit"

targets:
  - "antigravity"
  - "cursor"

architecture:
  pattern: "vertical-slice"
  strict_boundaries: true

workflows:
  - "sdd"
  - "tdd-strict"
  - "atomic-commits"

conventions:
  principles:
    - "solid-strict"
    - "immutability"
  naming: "pascal-case"

memory:
  provider: "sqlite-fts5"
  path: ".harness/memory.db"

guardrails:
  block_destructive_commands: true
  max_diff_lines: 300
  protected_paths:
    - ".git/**"
    - ".env*"
```

---

## 2. Motor de Plantillas y Embebido en Go

### 2.1 Embebido de Recursos (`embed.FS`)
Todas las plantillas maestras y fragmentos de reglas residen compilados dentro del binario final para garantizar portabilidad sin dependencias en el disco:

```go
package assets

import "embed"

//go:embed rules/*.tmpl skills/** templates/**
var AssetFS embed.FS
```

### 2.2 Motor de Renderizado (`engine.Renderer`)
Utiliza `text/template` de Go con funciones auxiliares personalizadas:
* `toLower`, `toUpper`, `camelCase`, `pascalCase`, `snakeCase`.
* Inclusión condicional de bloques basada en los flags de `harnessup.yaml`.

---

## 3. Emisores por Destino (Emitters)

Cada agente de IA tiene su propio adaptador emisor que implementa la interfaz `Emitter`:

```go
type Emitter interface {
    ID() string
    Emit(config *Config, renderer *Renderer, targetDir string) error
}
```

### 3.1 Adaptador Antigravity (`internal/emitter/antigravity.go`)
* Escribe reglas en `.agent/rules/`:
  * `.agent/rules/01-architecture.md`
  * `.agent/rules/02-workflow.md`
  * `.agent/rules/03-guardrails.md`
* Instala las carpetas de skills en `skills/` (cada una con su `SKILL.md`).
* Genera la configuración de Mission Control en `.agent/config.json`.

### 3.2 Adaptador Cursor (`internal/emitter/cursor.go`)
* Escribe reglas en `.cursor/rules/`:
  * Genera archivos `.mdc` con *frontmatter* que incluye `globs` específicos:
    ```markdown
    ---
    description: Reglas de Arquitectura Vertical Slice
    globs: src/Features/**/*.cs
    alwaysApply: true
    ---
    ```
* Escribe archivo legacy `.cursorrules` para máxima compatibilidad con versiones anteriores.

### 3.3 Adaptador Claude Code (`internal/emitter/claude.go`)
* Genera el archivo central `CLAUDE.md` con las instrucciones obligatorias y comandos de test y build.
* Genera subagentes y skills en `.claude/skills/`.

---

## 4. Estrategia de Idempotencia y Fusión (File Merge)

Para permitir que el usuario personalice sus reglas sin que `harnessup sync` borre sus cambios:
1. **Zonas protegidas por etiquetas:** Los archivos generados contendrán delimitadores:
   ```markdown
   <!-- HARNESSUP:MANAGED_START -->
   ... Reglas administradas automáticamente por HarnessUp ...
   <!-- HARNESSUP:MANAGED_END -->

   <!-- HARNESSUP:CUSTOM_USER_RULES -->
   ... Cualquier regla que el usuario agregue aquí nunca será sobreescrita ...
   ```
2. Si un archivo ya existe y no tiene etiquetas, se genera una copia de respaldo `.bak` antes de fusionar.
