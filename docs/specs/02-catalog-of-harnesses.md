# SPEC-002: Catálogo Exhaustivo de Arneses y Taxonomía

* **Proyecto:** HarnessUp
* **Estado:** Propuesto (Pendiente de Aprobación)
* **Fecha:** 2026-10-14
* **Metodología:** Spec-Driven Development (SDD)

---

## 1. Definición y Propósito del Catálogo

El catálogo de **HarnessUp** define los bloques modulares que el usuario puede combinar para armar el arnés de su proyecto. Cada bloque se compone de:
* Un identificador único (`ID`).
* Nombre legible y descripción.
* Inyecciones de reglas (prompts/guidelines para el LLM).
* Skills asociadas (`SKILL.md` autoportantes).
* Verificadores deterministas (scripts, linters, comandos de test o hooks).

---

## 2. Taxonomía de Arneses

### Categoría A: Arneses de Memoria Persistente (`memory`)

| ID | Nombre | Descripción | Implementación Técnica |
| :--- | :--- | :--- | :--- |
| `memory-sqlite-fts5` | SQLite FTS5 (Local-First) | Base de datos local SQLite con índice de búsqueda de texto completo FTS5 para guardar decisiones, bitácora y lecciones aprendidas. | Archivo `.harness/memory.db` gestionado por HarnessUp + tool CLI de consulta rápida. |
| `memory-markdown-flat` | Flat-File Markdown | Sistema liviano de archivos Markdown estructurados (`docs/decisions/` y `docs/memory.md`). | Plantillas preconfiguradas y reglas de actualización obligatoria post-tarea. |
| `memory-none` | Sin memoria persistente | Confía únicamente en el contexto efímero de la sesión actual. | Sin persistencia. |

---

### Categoría B: Arneses de Metodología y Flujo de Trabajo (`workflow`)

| ID | Nombre | Principio Rector | Reglas y Skills Generadas |
| :--- | :--- | :--- | :--- |
| `workflow-tdd-strict` | TDD Estricto (Red $\rightarrow$ Green $\rightarrow$ Refactor) | Prohibido escribir código de producción sin un test fallando primero. | Skill `tdd-workflow/SKILL.md`. Bloquea la aprobación sin evidencia de test en rojo previo. |
| `workflow-odd` | Organic-Driven Development (ODD) | Tareas pequeñas inline; tareas sustanciales con documento de feature, trazabilidad y commits unitarios. | Skill `odd-workflow/SKILL.md` + plantilla `odd/tasks/<feature>.md`. |
| `workflow-sdd` | Spec-Driven Development (SDD) | Prohibido codear sin especificación técnica validada por el humano. | Skill `sdd-workflow/SKILL.md` + carpeta `specs/` y reglas de aprobación previa. |
| `workflow-rdd-bounded` | Receipt-Driven Development (RDD) | Exige un "recibo" (ejecución de tests + diff acotado < 300 líneas) para cada unidad de trabajo. | Skill `work-unit-commits/SKILL.md` + verificador de tamaño de diff. |

---

### Categoría C: Arneses de Patrones de Arquitectura (`architecture`)

| ID | Nombre | Lenguajes Sugeridos | Restricciones Inyectadas al Agente |
| :--- | :--- | :--- | :--- |
| `arch-vertical-slice` | Vertical Slice Architecture | C#, TypeScript, Go | Todo código se agrupa por Feature (Command/Query, Handler, DTO, Endpoint, Tests). Prohibido crear carpetas genéricas `services/`, `repositories/` o `helpers/`. |
| `arch-clean` | Clean Architecture (Onion) | C#, Java, TypeScript | 4 capas estrictas: Domain $\rightarrow$ Application $\rightarrow$ Infrastructure $\rightarrow$ Presentation. La capa de Dominio tiene cero dependencias externas. |
| `arch-hexagonal` | Hexagonal (Ports & Adapters) | Go, Python, Java | Núcleo con lógica de negocio y puertos (interfaces); adaptadores para DBs, APIs y colas. |
| `arch-modular-monolith` | Monolito Modular | Todos | Módulos desacoplados que solo se comunican por contratos públicos definidos. |

---

### Categoría D: Arneses de Principios y Convenciones de Código (`conventions`)

| ID | Nombre | Reglas Inyectadas |
| :--- | :--- | :--- |
| `conv-solid-strict` | Principios SOLID Estrictos | Responsabilidad única, Open/Closed, Sustitución de Liskov, Segregación de interfaces e Inversión de dependencias. Prohíbe clases dios o métodos de más de 30 líneas. |
| `conv-naming-standard` | Convenciones de Nomenclatura | Configurable por lenguaje: `PascalCase` para C#, `camelCase` para TS/JS, `snake_case` para Python/Rust, `kebab-case` para archivos. |
| `conv-immutability` | Inmutabilidad & Funcional | Favorece estructuras inmutables, `records`, funciones puras y evita efectos secundarios no declarados. |

---

### Categoría E: Arneses de Seguridad y Guardrails (`guardrails`)

| ID | Nombre | Mecanismo |
| :--- | :--- | :--- |
| `guard-destructive-block` | Bloqueo de Comandos Destructivos | Lista negra estricta inyectada en las reglas: prohibido `rm -rf`, `git reset --hard` sin confirmación, `drop database`, o editar archivos `.env`. |
| `guard-protected-paths` | Rutas Protegidas | Restringe a la IA para que solo modifique archivos dentro del scope acordado (ej. no tocar CI/CD o infraestructura si es un fix de UI). |
| `guard-precommit-verify` | Verificación Pre-Commit | Hook local que ejecuta `linter` y `test` automáticamente antes de permitir un commit. |

---

### Categoría F: Arneses de Clientes de IA / Destinos (`targets`)

| ID | Cliente / Herramienta | Archivos que emite HarnessUp |
| :--- | :--- | :--- |
| `target-antigravity` | Google Antigravity | `.agent/rules/` + `skills/` + `.agent/config.json` |
| `target-cursor` | Cursor IDE | `.cursor/rules/*.mdc` (con globs) + `.cursorrules` |
| `target-claude-code` | Claude Code (Anthropic) | `CLAUDE.md` + `.claude/skills/` |
| `target-opencode` | OpenCode / Kilo Code | `.opencode/` + `opencode.json` + `skills/` |
| `target-windsurf` | Windsurf | `.windsurfrules` + `skills/` |
| `target-vscode-copilot` | VS Code Copilot | `.github/copilot-instructions.md` |
| `target-pi` | Pi / Gentle Shell | Instrucciones nativas compatibles y extensiones |
| `target-all` | Multi-Agente Universal | Emite la configuración para todos los agentes simultáneamente |
