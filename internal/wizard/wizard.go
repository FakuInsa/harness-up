package wizard

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/detector"
)

// Run inicia la ejecución interactiva del wizard en la terminal y retorna el Config resultante.
func Run(dirPath string) (*config.Config, error) {
	return RunWithOptions(dirPath)
}

// RunWithOptions permite inicializar el programa con opciones adicionales (ej. tea.WithInput o tea.WithoutRenderer).
func RunWithOptions(dirPath string, opts ...tea.ProgramOption) (*config.Config, error) {
	d := detector.New()
	stack, _ := d.DetectDir(dirPath)

	projectName := filepath.Base(dirPath)
	if projectName == "." || projectName == "" || projectName == "/" || projectName == "\\" {
		projectName = "HarnessApp"
	}

	model := NewModel(projectName, stack)
	p := tea.NewProgram(model, opts...)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("wizard execution failed: %w", err)
	}

	wm, ok := finalModel.(*WizardModel)
	if !ok || wm.IsAborted() {
		return nil, fmt.Errorf("wizard was cancelled by user")
	}

	return wm.Config(), nil
}

