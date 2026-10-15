package emitter

import (
	"fmt"

	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/engine"
)

const (
	// ManagedStart marca el inicio del bloque gestionado por HarnessUp.
	ManagedStart = "<!-- HARNESSUP:MANAGED_START -->"
	// ManagedEnd marca el final del bloque gestionado por HarnessUp.
	ManagedEnd = "<!-- HARNESSUP:MANAGED_END -->"
	// CustomUserRulesTag marca el inicio de la zona protegida para reglas del usuario.
	CustomUserRulesTag = "<!-- HARNESSUP:CUSTOM_USER_RULES -->"
)

// Emitter define el contrato para cualquier adaptador que genera artefactos de un agente de IA.
type Emitter interface {
	// ID retorna el identificador único del agente o destino (ej. "antigravity", "cursor").
	ID() string
	// Emit genera físicamente los archivos del arnés en el directorio destino según la configuración.
	Emit(cfg *config.Config, r *engine.Renderer, targetDir string) error
}

// Registry gestiona la colección de emisores disponibles en HarnessUp.
type Registry struct {
	emitters map[string]Emitter
}

// NewDefaultRegistry construye un registro con los emisores oficiales configurados.
func NewDefaultRegistry() *Registry {
	reg := &Registry{
		emitters: make(map[string]Emitter),
	}
	reg.Register(NewAntigravityEmitter())
	reg.Register(NewCursorEmitter())
	return reg
}

// Register asocia un emisor al registro.
func (reg *Registry) Register(e Emitter) {
	reg.emitters[e.ID()] = e
}

// Get obtiene un emisor por su ID.
func (reg *Registry) Get(id string) (Emitter, bool) {
	e, exists := reg.emitters[id]
	return e, exists
}

// EmitTargets ejecuta la emisión para todos los targets definidos en la configuración.
func (reg *Registry) EmitTargets(cfg *config.Config, r *engine.Renderer, targetDir string) error {
	for _, target := range cfg.Targets {
		e, exists := reg.Get(target)
		if !exists {
			return fmt.Errorf("unknown target emitter: %q", target)
		}
		if err := e.Emit(cfg, r, targetDir); err != nil {
			return fmt.Errorf("failed to emit target %q: %w", target, err)
		}
	}
	return nil
}
