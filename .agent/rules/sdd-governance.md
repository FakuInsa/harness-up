# Arnés de Gobernanza del Proyecto: SDD y Dogfooding

Este documento y el [Way of Working (WoW)](../../docs/WAY_OF_WORKING.md) gobiernan el ciclo de desarrollo de **HarnessUp**.
Toda IA y desarrollador que opere en este repositorio debe cumplir estrictamente estas reglas.

---

## 1. Regla Suprema: Spec-Driven Development (SDD)
* **Prohibido escribir código fuente sin especificación aprobada:** Ningún archivo `.go` ni script de implementación puede crearse hasta que la especificación correspondiente en `docs/specs/` esté completa, revisada y aprobada por el usuario.
* **Trazabilidad:** Cada feature, función o componente implementado debe mapear directamente a un requerimiento de la especificación técnica.

---

## 2. Registro de Decisiones de Arquitectura (ADR)
* Toda decisión técnica estructural (lenguaje, librerías clave, estructura de base de datos, protocolos de exportación) debe registrarse inmediatamente en `docs/specs/DECISIONS.md` con su contexto, opciones consideradas y justificación.

---

## 3. Disciplina de Código: TDD y Evidencia (RDD)
* **Test-First (Rojo -> Verde -> Refactor):**
  1. Escribir el test unitario en Go (`*_test.go`).
  2. Ejecutar y verificar que falla (RED).
  3. Escribir el código mínimo para que pase (GREEN).
  4. Refactorizar manteniendo el test en verde.
* **Evidencia determinista:** Ninguna tarea se da por cerrada sin ejecutar los tests (`go test ./...`) y mostrar el resultado en consola.

---

## 4. Unidades de Trabajo Atómicas y Ramas (Branches)
* Branch por Feature: Cada feature debe implementarse en su propia rama Git (`feat/<id>-<nombre>`), la cual se fusiona a `main` únicamente tras la aprobación final de la evidencia.
* Los cambios deben ser acotados, pequeños y coherentes (< 300 líneas de diff productivo).
* No realizar refactorizaciones no solicitadas ni tocar módulos fuera del alcance de la tarea activa.

---

## 5. Prevención de Commit Bloat en Metadatos (ADR-005)
* **Prohibidos los micro-commits de metadatos:** No crear commits dedicados exclusivamente a modificar marcas de avance o estado de specs (`[ ]`, `[⏳]`, `[✓]`, o marcas de Gate).
* **Estructura atómica:** Máximo 2 commits por feature: (1) `docs(specs): add specification for feat-xxx` y (2) `feat(xxx): implement <feature>` integrando código, tests y cierre de estado.

