package wizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/detector"
)

// WizardModel encapsula la máquina de estados de terminal para configurar HarnessUp.
type WizardModel struct {
	step              Step
	stack             *detector.StackInfo
	config            *config.Config
	cursor            int
	selectedTargets   map[string]bool
	selectedWorkflows map[string]bool
	done              bool
	aborted           bool
}

// NewModel inicializa un nuevo WizardModel evaluando si existe stack previo (Brownfield vs Greenfield).
func NewModel(projectName string, stack *detector.StackInfo) *WizardModel {
	cfg := config.Default(projectName)

	initialStep := StepSelectLanguage
	if stack != nil && stack.Language != "" && stack.Language != detector.LangUnknown {
		cfg.Project.Stack.Language = string(stack.Language)
		cfg.Project.Stack.Framework = stack.Framework
		cfg.Project.Stack.TestRunner = string(stack.TestRunner)
		initialStep = StepDetectConfirm
	}

	targetsMap := make(map[string]bool)
	for _, t := range cfg.Targets {
		targetsMap[t] = true
	}

	workflowsMap := make(map[string]bool)
	for _, w := range cfg.Workflows {
		workflowsMap[w] = true
	}

	return &WizardModel{
		step:              initialStep,
		stack:             stack,
		config:            cfg,
		cursor:            0,
		selectedTargets:   targetsMap,
		selectedWorkflows: workflowsMap,
	}
}

// Step retorna el paso activo.
func (m *WizardModel) Step() Step { return m.step }

// Config retorna la configuración resultante construida.
func (m *WizardModel) Config() *config.Config { return m.config }

// IsDone indica si el asistente se completó satisfactoriamente.
func (m *WizardModel) IsDone() bool { return m.done }

// IsAborted indica si el usuario canceló la ejecución.
func (m *WizardModel) IsAborted() bool { return m.aborted }

// Init implementa tea.Model.
func (m *WizardModel) Init() tea.Cmd { return nil }

// UpdateModel es un helper tipado para testing headless.
func (m *WizardModel) UpdateModel(msg tea.Msg) *WizardModel {
	newM, _ := m.Update(msg)
	return newM.(*WizardModel)
}

// Update procesa los eventos reactivos de teclado.
func (m *WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case tea.KeyDown:
			if m.cursor < m.maxCursor() {
				m.cursor++
			}
			return m, nil
		case tea.KeySpace:
			m.toggleItem()
			return m, nil
		case tea.KeyEnter:
			return m.advance()
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				switch msg.Runes[0] {
				case 'k':
					if m.cursor > 0 {
						m.cursor--
					}
					return m, nil
				case 'j':
					if m.cursor < m.maxCursor() {
						m.cursor++
					}
					return m, nil
				case ' ':
					m.toggleItem()
					return m, nil
				case 'q':
					if m.step == StepSummary {
						m.aborted = true
						return m, tea.Quit
					}
				}
			}
		}
	}
	return m, nil
}

func (m *WizardModel) maxCursor() int {
	switch m.step {
	case StepDetectConfirm:
		return 1 // 0: Si usar detectado, 1: No, configurar manualmente
	case StepSelectLanguage:
		return len(languageOptions) - 1
	case StepSelectFramework:
		fwList := frameworksByLang[m.config.Project.Stack.Language]
		if len(fwList) == 0 {
			return 0
		}
		return len(fwList) - 1
	case StepSelectTargets:
		return len(targetOptions) - 1
	case StepSelectArchitecture:
		return len(architectureOptions) - 1
	case StepSelectWorkflows:
		return len(workflowOptions) - 1
	case StepSelectNaming:
		return len(namingOptions) - 1
	case StepSelectMemory:
		return len(memoryOptions) - 1
	case StepSummary:
		return 0
	}
	return 0
}

func (m *WizardModel) toggleItem() {
	switch m.step {
	case StepSelectTargets:
		if m.cursor < len(targetOptions) {
			key := targetOptions[m.cursor].Key
			m.selectedTargets[key] = !m.selectedTargets[key]
		}
	case StepSelectWorkflows:
		if m.cursor < len(workflowOptions) {
			key := workflowOptions[m.cursor].Key
			m.selectedWorkflows[key] = !m.selectedWorkflows[key]
		}
	}
}

func (m *WizardModel) advance() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepDetectConfirm:
		if m.cursor == 0 {
			m.step = StepSelectTargets
		} else {
			m.step = StepSelectLanguage
		}
	case StepSelectLanguage:
		if m.cursor < len(languageOptions) {
			m.config.Project.Stack.Language = languageOptions[m.cursor].Key
		}
		m.step = StepSelectFramework
	case StepSelectFramework:
		fwList := frameworksByLang[m.config.Project.Stack.Language]
		if m.cursor < len(fwList) {
			m.config.Project.Stack.Framework = fwList[m.cursor].Key
		}
		m.step = StepSelectTargets
	case StepSelectTargets:
		var targets []string
		for _, opt := range targetOptions {
			if m.selectedTargets[opt.Key] {
				targets = append(targets, opt.Key)
			}
		}
		if len(targets) == 0 {
			targets = []string{"antigravity"}
		}
		m.config.Targets = targets
		m.step = StepSelectArchitecture
	case StepSelectArchitecture:
		if m.cursor < len(architectureOptions) {
			m.config.Architecture.Pattern = architectureOptions[m.cursor].Key
		}
		m.step = StepSelectWorkflows
	case StepSelectWorkflows:
		var workflows []string
		for _, opt := range workflowOptions {
			if m.selectedWorkflows[opt.Key] {
				workflows = append(workflows, opt.Key)
			}
		}
		if len(workflows) == 0 {
			workflows = []string{"sdd"}
		}
		m.config.Workflows = workflows
		m.step = StepSelectNaming
	case StepSelectNaming:
		if m.cursor < len(namingOptions) {
			m.config.Conventions.Naming = namingOptions[m.cursor].Key
		}
		m.step = StepSelectMemory
	case StepSelectMemory:
		if m.cursor < len(memoryOptions) {
			m.config.Memory.Provider = memoryOptions[m.cursor].Key
		}
		m.step = StepSummary
	case StepSummary:
		m.done = true
		m.step = StepFinished
		return m, tea.Quit
	}

	m.cursor = 0
	return m, nil
}
