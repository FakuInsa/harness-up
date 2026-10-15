# Architectural Decision Records (ADR) - HarnessUp

Este documento registra cronológicamente las decisiones técnicas y de arquitectura fundamentales del proyecto **HarnessUp**.

---

## ADR-001: Adopción Estricta de Spec-Driven Development (SDD) y Dogfooding

* **Fecha:** 2026-10-14
* **Estado:** Aceptado
* **Contexto:**
  Para evitar la deriva del proyecto, el código no planificado y la sobreingeniería típica al trabajar con agentes de IA, se requiere una metodología que garantice rigor, trazabilidad y control.
* **Decisión:**
  Se adopta **Spec-Driven Development (SDD)** como metodología rectora. No se escribirá una sola línea de código en Go hasta que las especificaciones funcionales, de arquitectura, de catálogo y de CLI estén completamente redactadas y aprobadas en `docs/specs/`. Asimismo, el proyecto aplica *dogfooding*: se gobierna a sí mismo utilizando un arnés de reglas en `.agent/rules/sdd-governance.md`.
* **Consecuencias:**
  * Mayor claridad arquitectónica y reducción drástica de refactorizaciones futuras.
  * Trazabilidad completa entre requerimientos, tests y código de producción.

---

## ADR-002: Elección del Stack Tecnológico: Go + Bubbletea + SQLite FTS5

* **Fecha:** 2026-10-14
* **Estado:** Aceptado
* **Contexto:**
  La herramienta debe ser una CLI interactiva, rápida, multiplataforma, estética y capaz de ejecutarse como un binario único sin requerir runtimes pesados ni dependencias en el sistema del usuario.
* **Opciones Evaluadas:**
  1. *TypeScript / Node.js:* Excelente por `npx`, pero requiere Node.js instalado y dependencia de `node_modules`.
  2. *C# (.NET 9) con Spectre.Console:* Excelente calidad visual y Native AOT, pero ecosistema menos natural para usuarios fuera de .NET.
  3. *Go (Golang) con Bubbletea (Charm.sh):* Compila a binario único nativo, arranque instantáneo (0 ms), excelente ecosistema de TUI (usado por Gentle-AI y herramientas de referencia), y soporte embebido de SQLite.
  4. *Python:* Descartado por complejidad de distribución y peso de binarios.
* **Decisión:**
  Se elige **Go** como lenguaje principal, con el siguiente ecosistema central:
  * **TUI / Wizard:** `github.com/charmbracelet/bubbletea` y `github.com/charmbracelet/lipgloss` para interfaces de terminal enriquecidas y accesibles.
  * **CLI Framework:** `github.com/spf13/cobra` para comandos, flags y subcomandos (`init`, `sync`, `doctor`, `list`).
  * **Persistencia / Memoria Local:** SQLite nativo embebido con FTS5 (Full-Text Search 5) para indexación de memoria y decisiones.
  * **Plantillas:** `text/template` nativo de Go y `embed.FS` para empaquetar todas las skills, reglas y configuraciones dentro del binario compilado.
* **Consecuencias:**
  * Binario único, distribuible vía GitHub Releases, Homebrew, `go install`, o descarga directa de `.exe`.
  * Rendimiento extremo en arranque y ejecución de hooks locales.

---

## ADR-003: Estrategia de Ramas Git por Feature (Branch-per-Feature)

* **Fecha:** 2026-10-15
* **Estado:** Aceptado
* **Contexto:**
  Para mantener la rama `main` siempre estable y facilitar la trazabilidad estricta de cada funcionalidad del roadmap en el ciclo SDD/TDD, se necesita una política clara de aislamiento de código.
* **Decisión:**
  Se adopta el patrón **Branch-per-Feature**:
  * Cada feature se desarrolla en su propia rama `feat/<id>-<nombre-kebab>` (ej. `feat/feat-001-stack-detector`).
  * La rama se crea tras la aprobación de la spec (Fase 2: GATE).
  * Todos los ciclos RED/GREEN y commits de la unidad de trabajo viven en la rama.
  * La integración a `main` ocurre únicamente tras la aprobación final de la evidencia (Fase 5) por parte del Tech Lead.
* **Consecuencias:**
  * La rama `main` siempre permanece en un estado 100% verde y desplegable.
  * Historial limpio, auditable e integrable mediante PR o merge directo.

---

## ADR-004: Soporte Dual Brownfield / Greenfield en el Asistente de Configuración

* **Fecha:** 2026-10-15
* **Estado:** Aceptado
* **Contexto:**
  HarnessUp debe operar eficazmente tanto en repositorios existentes con código fuente previo (**Modo Brownfield**) como en directorios vacíos o proyectos nuevos que se inician desde cero (**Modo Greenfield**). En este último caso, el arnés de gobernanza, reglas de arquitectura, skills y convenciones de nomenclatura deben configurarse *antes* de que el agente de IA empiece a escribir código.
* **Decisión:**
  1. Si el escáner del detector reporta stack `unknown` o si el usuario invoca `harnessup init --new`, la TUI no aborta ni deja campos vacíos; en su lugar, despliega un flujo asistido de preguntas (Stack, Framework, Arquitectura estándar o híbrida, y Nomenclatura como snake_case).
  2. Se desacopla la generación del **arnés de reglas de IA** (alcance núcleo de HarnessUp v1) del **scaffolding físico de código y plantillas del framework** (ej. ejecutar `dotnet new`), difiriendo este último a una feature futura de extensión (`FEAT-008`).
* **Consecuencias:**
  * Máxima flexibilidad: los proyectos nacen gobernados desde el minuto cero.
  * Claridad en los límites del arnés sin acoplar la v1 a instaladores pesados de SDKs de terceros.

---

## ADR-005: Prevención de Commit Bloat en Ciclos SDD (Granularidad Atómica)

* **Fecha:** 2026-10-15
* **Estado:** Aceptado
* **Contexto:**
  Durante las fases iniciales de desarrollo (FEAT-001 a FEAT-005), la aplicación ultra-literal de cada fase del ciclo SDD (SPEC -> GATE -> RED -> GREEN -> DONE) generó micro-commits aislados dedicados únicamente a cambiar estados en markdown (ej. `[ ]` a `[⏳]` o `[✓]`). Esto produjo ruido en el árbol de Git (*commit bloat*) con commits carentes de valor funcional independiente.
* **Decisión:**
  1. Se prohíbe terminantemente la creación de micro-commits dedicados exclusivamente a actualizar metadatos o estados de avance (`docs(sdd): approve spec`, `docs(roadmap): mark done`).
  2. Cada unidad de trabajo o feature se compone como máximo de **2 commits atómicos**:
     * **Commit 1 (Especificación):** `docs(specs): add specification for feat-xxx` (redacta la especificación técnica en `docs/specs/`).
     * **Commit 2 (Implementación Completa):** `feat(xxx): implement <feature>` (incluye tests unitarios, código de producción, y la actualización a `DONE` de la especificación y del roadmap).
  3. Los arneses y plantillas generadas por HarnessUp deben incorporar explícitamente esta regla en sus guardarraíles para evitar que agentes downstream repitan esta práctica.
* **Consecuencias:**
  * Historial de Git limpio, semántico y legible donde cada commit representa una entrega funcional verificable.
  * Eliminación total del ruido de commits de metadatos administrativos.

