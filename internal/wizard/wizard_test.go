package wizard_test

import (
	"bytes"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/FakuInsa/harness-up/internal/detector"
	"github.com/FakuInsa/harness-up/internal/wizard"
)

func sendKey(m tea.Model, msg tea.KeyMsg) tea.Model {
	updated, _ := m.Update(msg)
	return updated
}

func sendEnter(m tea.Model) tea.Model {
	return sendKey(m, tea.KeyMsg{Type: tea.KeyEnter})
}

func sendSpace(m tea.Model) tea.Model {
	return sendKey(m, tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})
}

func sendDown(m tea.Model) tea.Model {
	return sendKey(m, tea.KeyMsg{Type: tea.KeyDown})
}

func sendUp(m tea.Model) tea.Model {
	return sendKey(m, tea.KeyMsg{Type: tea.KeyUp})
}

func TestWizard_BrownfieldFlow(t *testing.T) {
	stack := &detector.StackInfo{
		Language:   detector.LangCSharp,
		Framework:  "dotnet-9",
		TestRunner: detector.RunnerXUnit,
		Confidence: 0.95,
	}

	m := wizard.NewModel("BrownfieldApp", stack)

	if m.Step() != wizard.StepDetectConfirm {
		t.Fatalf("expected initial step to be StepDetectConfirm, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepDetectConfirm")
	}

	if m.Init() != nil {
		t.Error("expected nil init cmd")
	}

	// 1. Confirm detected stack (Enter) -> Targets
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectTargets {
		t.Fatalf("expected StepSelectTargets, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSelectTargets")
	}

	// 2. Accept default targets (Enter) -> Architecture
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectArchitecture {
		t.Fatalf("expected StepSelectArchitecture, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSelectArchitecture")
	}

	// 3. Accept default architecture (Enter) -> Workflows
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectWorkflows {
		t.Fatalf("expected StepSelectWorkflows, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSelectWorkflows")
	}

	// 4. Accept default workflows (Enter) -> Naming
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectNaming {
		t.Fatalf("expected StepSelectNaming, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSelectNaming")
	}

	// 5. Accept default naming (Enter) -> Memory
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectMemory {
		t.Fatalf("expected StepSelectMemory, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSelectMemory")
	}

	// 6. Accept default memory (Enter) -> Summary
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSummary {
		t.Fatalf("expected StepSummary, got %v", m.Step())
	}
	if m.View() == "" {
		t.Error("expected non-empty view in StepSummary")
	}

	// 7. Confirm summary (Enter) -> Finished
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.IsDone() {
		t.Fatalf("expected wizard to be done")
	}
	if m.IsAborted() {
		t.Fatalf("wizard should not be aborted")
	}
	if m.View() == "" {
		t.Error("expected non-empty view when done")
	}

	cfg := m.Config()
	if cfg == nil {
		t.Fatalf("expected non-nil config")
	}
	if cfg.Project.Name != "BrownfieldApp" {
		t.Errorf("expected project name BrownfieldApp, got %q", cfg.Project.Name)
	}
	if cfg.Project.Stack.Language != "csharp" {
		t.Errorf("expected stack language csharp, got %q", cfg.Project.Stack.Language)
	}
	if cfg.Project.Stack.Framework != "dotnet-9" {
		t.Errorf("expected stack framework dotnet-9, got %q", cfg.Project.Stack.Framework)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("resulting config should be valid, got: %v", err)
	}
}

func TestWizard_GreenfieldFlow(t *testing.T) {
	stack := &detector.StackInfo{
		Language: detector.LangUnknown,
	}

	m := wizard.NewModel("GreenfieldApp", stack)

	// Since stack is unknown, starts at StepSelectLanguage
	if m.Step() != wizard.StepSelectLanguage {
		t.Fatalf("expected initial step to be StepSelectLanguage, got %v", m.Step())
	}

	// Navigate to Go (csharp -> typescript -> python -> go = 3 down)
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Select Go -> Framework

	if m.Step() != wizard.StepSelectFramework {
		t.Fatalf("expected StepSelectFramework, got %v", m.Step())
	}

	// Accept default framework (standard) -> Targets
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	if m.Step() != wizard.StepSelectTargets {
		t.Fatalf("expected StepSelectTargets, got %v", m.Step())
	}

	// Advance through remaining steps
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Arch
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Workflows
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Naming
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Memory
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Summary
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Done

	if !m.IsDone() {
		t.Fatalf("expected wizard to be done")
	}

	cfg := m.Config()
	if cfg == nil {
		t.Fatalf("expected non-nil config")
	}
	if cfg.Project.Stack.Language != "go" {
		t.Errorf("expected language go, got %q", cfg.Project.Stack.Language)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("resulting config should be valid, got: %v", err)
	}
}

func TestWizard_Brownfield_RejectDetected(t *testing.T) {
	stack := &detector.StackInfo{
		Language:   detector.LangPython,
		Framework:  "django",
		TestRunner: detector.RunnerPytest,
		Confidence: 0.9,
	}

	m := wizard.NewModel("RejectApp", stack)
	if m.Step() != wizard.StepDetectConfirm {
		t.Fatalf("expected StepDetectConfirm, got %v", m.Step())
	}

	// Move cursor to "No, configure manually" (Down) and Enter
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})

	// Should transition to StepSelectLanguage
	if m.Step() != wizard.StepSelectLanguage {
		t.Fatalf("expected StepSelectLanguage after rejecting detection, got %v", m.Step())
	}
}

func TestWizard_MultiselectToggle(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("ToggleApp", stack)

	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Lang -> Fw
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Fw -> Targets

	if m.Step() != wizard.StepSelectTargets {
		t.Fatalf("expected StepSelectTargets, got %v", m.Step())
	}

	// Move cursor down to "cursor" (target 2) and press Space to toggle it on
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})

	// Advance through remaining steps
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Arch
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Workflows
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Naming
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Memory
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Summary
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Done

	cfg := m.Config()
	if cfg == nil {
		t.Fatalf("expected non-nil config")
	}

	hasCursor := false
	for _, tgt := range cfg.Targets {
		if tgt == "cursor" {
			hasCursor = true
			break
		}
	}
	if !hasCursor {
		t.Errorf("expected targets to contain 'cursor' after toggling, got %v", cfg.Targets)
	}
}

func TestWizard_Abort(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("AbortApp", stack)

	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !m.IsAborted() {
		t.Errorf("expected wizard to be aborted on ctrl+c")
	}
	if m.IsDone() {
		t.Errorf("wizard should not be done when aborted")
	}
}

func TestWizard_NavigationBounds(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("NavApp", stack)

	// Move up at top shouldn't panic or go below 0
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyUp})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyUp})

	// Move down multiple times with 'j' and 'k'
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	view := m.View()
	if view == "" {
		t.Errorf("expected non-empty view string")
	}
}

func TestWizard_SummaryAbortWithQ(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("QApp", stack)
	for i := 0; i < 7; i++ {
		m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter})
	}
	if m.Step() != wizard.StepSummary {
		t.Fatalf("expected StepSummary, got %v", m.Step())
	}
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if !m.IsAborted() {
		t.Errorf("expected aborted on q in summary")
	}
}

func TestWizard_ToggleWorkflows(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("WorkflowApp", stack)
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Lang -> Fw
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Fw -> Targets
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Targets -> Arch
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Arch -> Workflows

	if m.Step() != wizard.StepSelectWorkflows {
		t.Fatalf("expected StepSelectWorkflows, got %v", m.Step())
	}
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyDown})
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}})

	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Naming
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Memory
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Summary
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEnter}) // Done

	cfg := m.Config()
	if len(cfg.Workflows) == 0 {
		t.Errorf("expected fallback workflow when none selected")
	}
}

func TestWizard_AbortWithEsc(t *testing.T) {
	stack := &detector.StackInfo{Language: detector.LangUnknown}
	m := wizard.NewModel("EscApp", stack)
	m = m.UpdateModel(tea.KeyMsg{Type: tea.KeyEsc})
	if !m.IsAborted() {
		t.Errorf("expected aborted on esc")
	}
	if m.View() == "" {
		t.Errorf("expected view on abort")
	}
}

func TestRunWithOptions_Success(t *testing.T) {
	tmpDir := t.TempDir()
	// Provide enters to navigate through the entire greenfield wizard
	input := bytes.NewBufferString("\r\r\r\r\r\r\r\r\r")
	cfg, err := wizard.RunWithOptions(tmpDir, tea.WithInput(input), tea.WithoutRenderer())
	if err != nil {
		t.Fatalf("unexpected error running wizard: %v", err)
	}
	if cfg == nil {
		t.Fatalf("expected non-nil config")
	}
}

func TestRunWithOptions_Abort(t *testing.T) {
	tmpDir := t.TempDir()
	input := bytes.NewBufferString("\x03") // Ctrl+C
	_, err := wizard.RunWithOptions(tmpDir, tea.WithInput(input), tea.WithoutRenderer())
	if err == nil {
		t.Fatalf("expected error on abort")
	}
}


