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
