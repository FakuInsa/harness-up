# Tablero de Features (Kanban / Trello) — HarnessUp

Este archivo funciona como nuestro **tablero Trello**. Cada feature debe pasar secuencialmente por el ciclo SDD:
**[ ] TO DO** $\rightarrow$ **[⏳] IN PROGRESS (RED -> GREEN)** $\rightarrow$ **[✓] DONE (Verificado con Evidencia)**.

No se inicia una feature hasta que la anterior esté en **[✓] DONE** y aprobada por el Tech Lead.

---

## Columna: Backlog de Features

### Feature 1: Scaffolding Base & Detector de Stack
* **ID:** `FEAT-001`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Inicializar el módulo Go y construir el componente `internal/detector` que escanea el directorio para reconocer el lenguaje (C#, TypeScript, Python, Go) y test runner existente.
* **Criterio de Aceptación:**
  * Tests unitarios (`detector_test.go`) con fixtures simulados de cada stack.
  * Pasa la prueba al 100% con `go test ./internal/detector/...`.

---

### Feature 2: Modelo de Configuración & Serialización (`harnessforge.yaml`)
* **ID:** `FEAT-002`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Definir los structs de Go para la configuración del arnés y programar la lectura, escritura y validación del archivo `harnessforge.yaml`.
* **Criterio de Aceptación:**
  * Tests unitarios validando carga y guardado idempotente sin pérdida de datos.
  * Rechazo con error claro ante configuraciones inválidas.

---

### Feature 3: Motor de Plantillas y Assets Embebidos (`embed.FS`)
* **ID:** `FEAT-003`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Empaquetar las primeras plantillas de reglas (`rules/`) y skills (`skills/`) dentro del binario Go con `//go:embed` y crear el motor de compilación `internal/engine`.
* **Criterio de Aceptación:**
  * Renderizado exitoso de variables (`{{.Project.Name}}`, `{{.Architecture}}`).
  * Tests unitarios verificando que las plantillas embebidas se resuelven en memoria.

---

### Feature 4: Emisores de Archivos (Antigravity & Cursor)
* **ID:** `FEAT-004`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Programar los adaptadores emisores que toman la configuración y generan físicamente los archivos `.agent/rules/`, `.cursor/rules/*.mdc` y carpetas de skills con etiquetas de fusión protegidas.
* **Criterio de Aceptación:**
  * Tests verificando creación de archivos en directorio temporal.
  * Respeto a las zonas `<!-- HARNESSFORGE:CUSTOM_USER_RULES -->` sin sobreescribir código del usuario.

---

### Feature 5: Wizard Interactivo TUI (Bubbletea + Lipgloss)
* **ID:** `FEAT-005`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Construir la máquina de estados interactiva en consola (selección de agente, arquitectura, metodología y convenciones con flechas y checkboxes).
* **Criterio de Aceptación:**
  * Navegación fluida por teclado sin cuelgues.
  * Generación del objeto `Config` resultante al presionar Enter en el último paso.

---

### Feature 6: CLI & Comandos Cobra (`init`, `sync`, `doctor`)
* **ID:** `FEAT-006`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Conectar el Wizard y los Emisores al comando principal `harnessforge init` y añadir flags para modo no interactivo (`--non-interactive`).
* **Criterio de Aceptación:**
  * Ejecutable de consola probado de punta a punta.

---

### Feature 7: Módulo de Memoria Persistente (SQLite FTS5)
* **ID:** `FEAT-007`
* **Estado:** `[ ] TO DO`
* **Objetivo:** Crear la base de datos local `.harness/memory.db`, tabla FTS5 de decisiones, e interfaz CLI para guardar y consultar memoria técnica.
* **Criterio de Aceptación:**
  * Consultas de búsqueda de texto completo con tiempos de respuesta menores a 5 ms.
