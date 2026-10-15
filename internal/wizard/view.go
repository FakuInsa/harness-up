package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDD"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Italic(true)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

// View renderiza la interfaz visual del wizard según el paso activo.
func (m *WizardModel) View() string {
	if m.done {
		return headerStyle.Render("✓ ¡Configuración del arnés completada con éxito!\n")
	}
	if m.aborted {
		return descStyle.Render("Operación cancelada por el usuario.\n")
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("⚡ HarnessUp — Deterministic AI Engineering Harness"))
	b.WriteString("\n")

	switch m.step {
	case StepDetectConfirm:
		b.WriteString(headerStyle.Render("Paso 1: Detección de Proyecto Existente (Modo Brownfield)"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("• Lenguaje detectado: %s\n", m.config.Project.Stack.Language))
		if m.config.Project.Stack.Framework != "" {
			b.WriteString(fmt.Sprintf("• Framework detectado: %s\n", m.config.Project.Stack.Framework))
		}
		if m.config.Project.Stack.TestRunner != "" {
			b.WriteString(fmt.Sprintf("• Suite de pruebas: %s\n", m.config.Project.Stack.TestRunner))
		}
		b.WriteString("\n¿Deseas inicializar el arnés con esta configuración detectada?\n\n")

		options := []string{"Sí, usar esta configuración", "No, configurar manualmente (Greenfield)"}
		for i, opt := range options {
			cursor := "  "
			style := normalStyle
			if m.cursor == i {
				cursor = "❯ "
				style = selectedStyle
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(opt)))
		}

	case StepSelectLanguage:
		b.WriteString(headerStyle.Render("Paso 1B: Selección de Lenguaje (Modo Greenfield)"))
		b.WriteString("\n¿Qué ecosistema utilizará este proyecto?\n\n")
		m.renderSingleSelect(&b, languageOptions)

	case StepSelectFramework:
		b.WriteString(headerStyle.Render("Paso 1C: Selección de Framework Base"))
		b.WriteString("\n¿Qué framework o plantilla base gobierna el proyecto?\n\n")
		fwList := frameworksByLang[m.config.Project.Stack.Language]
		m.renderSingleSelect(&b, fwList)

	case StepSelectTargets:
		b.WriteString(headerStyle.Render("Paso 2: Selección de Agentes de IA"))
		b.WriteString("\n¿Qué agentes o IDEs asistidos gobernarán el proyecto? (Usa Espacio para seleccionar)\n\n")
		for i, opt := range targetOptions {
			cursor := "  "
			style := normalStyle
			if m.cursor == i {
				cursor = "❯ "
				style = selectedStyle
			}
			checked := "[ ]"
			if m.selectedTargets[opt.Key] {
				checked = "[✓]"
			}
			b.WriteString(fmt.Sprintf("%s%s %s — %s\n", cursor, checked, style.Render(opt.Label), descStyle.Render(opt.Description)))
		}

	case StepSelectArchitecture:
		b.WriteString(headerStyle.Render("Paso 3: Patrón de Arquitectura"))
		b.WriteString("\n¿Qué patrón estructural debe imponer el arnés al agente de IA?\n\n")
		m.renderSingleSelect(&b, architectureOptions)

	case StepSelectWorkflows:
		b.WriteString(headerStyle.Render("Paso 4: Metodología y Workflows de Desarrollo"))
		b.WriteString("\n¿Qué políticas metodológicas deben aplicarse estrictamente? (Usa Espacio para seleccionar)\n\n")
		for i, opt := range workflowOptions {
			cursor := "  "
			style := normalStyle
			if m.cursor == i {
				cursor = "❯ "
				style = selectedStyle
			}
			checked := "[ ]"
			if m.selectedWorkflows[opt.Key] {
				checked = "[✓]"
			}
			b.WriteString(fmt.Sprintf("%s%s %s — %s\n", cursor, checked, style.Render(opt.Label), descStyle.Render(opt.Description)))
		}

	case StepSelectNaming:
		b.WriteString(headerStyle.Render("Paso 5: Convenciones de Nomenclatura (Naming)"))
		b.WriteString("\n¿Qué estilo de nomenclatura debe exigir el arnés en código y archivos?\n\n")
		m.renderSingleSelect(&b, namingOptions)

	case StepSelectMemory:
		b.WriteString(headerStyle.Render("Paso 6: Persistencia y Memoria Técnica"))
		b.WriteString("\n¿Cómo se indexarán las decisiones de arquitectura y memoria del proyecto?\n\n")
		m.renderSingleSelect(&b, memoryOptions)

	case StepSummary:
		b.WriteString(headerStyle.Render("Paso 7: Resumen y Confirmación Final"))
		b.WriteString("\n")

		summaryText := fmt.Sprintf(
			"Proyecto:     %s\n"+
				"Lenguaje:     %s\n"+
				"Framework:    %s\n"+
				"Arquitectura: %s\n"+
				"Agentes:      %s\n"+
				"Workflows:    %s\n"+
				"Naming:       %s\n"+
				"Memoria:      %s",
			m.config.Project.Name,
			m.config.Project.Stack.Language,
			m.config.Project.Stack.Framework,
			m.config.Architecture.Pattern,
			strings.Join(m.config.Targets, ", "),
			strings.Join(m.config.Workflows, ", "),
			m.config.Conventions.Naming,
			m.config.Memory.Provider,
		)

		b.WriteString(cardStyle.Render(summaryText))
		b.WriteString("\n\nPresiona [Enter] para generar el arnés o [Esc/q] para cancelar.\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓/j/k: Navegar • Espacio: Alternar • Enter: Continuar • Esc/Ctrl+C: Salir"))
	return b.String()
}

func (m *WizardModel) renderSingleSelect(b *strings.Builder, items []OptionItem) {
	for i, opt := range items {
		cursor := "  "
		style := normalStyle
		bullet := "( )"
		if m.cursor == i {
			cursor = "❯ "
			style = selectedStyle
			bullet = "(•)"
		}
		b.WriteString(fmt.Sprintf("%s%s %s — %s\n", cursor, bullet, style.Render(opt.Label), descStyle.Render(opt.Description)))
	}
}
