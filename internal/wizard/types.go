package wizard

// Step identifica cada etapa secuencial del asistente.
type Step int

const (
	StepDetectConfirm Step = iota
	StepSelectLanguage
	StepSelectFramework
	StepSelectTargets
	StepSelectArchitecture
	StepSelectWorkflows
	StepSelectNaming
	StepSelectMemory
	StepSummary
	StepFinished
)

// OptionItem representa un elemento elegible en una lista.
type OptionItem struct {
	Key         string
	Label       string
	Description string
}

var languageOptions = []OptionItem{
	{Key: "csharp", Label: "C# / .NET", Description: "Ecosistema .NET moderno"},
	{Key: "typescript", Label: "TypeScript / Node.js", Description: "Backend & Fullstack JS/TS"},
	{Key: "python", Label: "Python", Description: "APIs y automatización con Python"},
	{Key: "go", Label: "Go (Golang)", Description: "Servicios concurrentes y CLIs de alto rendimiento"},
}

var frameworksByLang = map[string][]OptionItem{
	"csharp": {
		{Key: "dotnet-9", Label: "ASP.NET Core (.NET 9)", Description: "Web APIs y microservicios"},
		{Key: "minimal-apis", Label: "Minimal APIs", Description: "Servicios ligeros"},
		{Key: "console", Label: "Console Application", Description: "Herramientas de línea de comandos"},
	},
	"typescript": {
		{Key: "nestjs", Label: "NestJS", Description: "Framework empresarial TypeScript"},
		{Key: "express", Label: "Express / Fastify", Description: "Servidor HTTP minimalista"},
		{Key: "nextjs", Label: "Next.js", Description: "Fullstack React framework"},
	},
	"python": {
		{Key: "fastapi", Label: "FastAPI", Description: "APIs asíncronas modernas"},
		{Key: "django", Label: "Django", Description: "Framework web completo"},
		{Key: "flask", Label: "Flask", Description: "Microframework web flexible"},
	},
	"go": {
		{Key: "standard", Label: "Standard Library (net/http)", Description: "Sin dependencias externas"},
		{Key: "gin", Label: "Gin Web Framework", Description: "Enrutador HTTP rápido"},
		{Key: "fiber", Label: "Fiber", Description: "Framework inspirado en Express sobre Fasthttp"},
	},
}

var targetOptions = []OptionItem{
	{Key: "antigravity", Label: "Google Antigravity", Description: "Arnés .agent/rules/ y skills/"},
	{Key: "cursor", Label: "Cursor IDE", Description: "Reglas .cursor/rules/*.mdc con frontmatter"},
}

var architectureOptions = []OptionItem{
	{Key: "vertical-slice", Label: "Vertical Slice Architecture", Description: "Features autónomas (Commands, Queries, DTOs)"},
	{Key: "clean", Label: "Clean Architecture", Description: "Capas concéntricas (Domain, App, Infra)"},
	{Key: "hexagonal", Label: "Hexagonal (Ports & Adapters)", Description: "Núcleo aislado mediante interfaces"},
	{Key: "hybrid", Label: "Híbrida / Modular", Description: "Arquitectura combinada según el módulo"},
}

var workflowOptions = []OptionItem{
	{Key: "sdd", Label: "Spec-Driven Development", Description: "Prohibido codear sin spec técnica aprobada"},
	{Key: "tdd-strict", Label: "TDD Estricto", Description: "Red -> Green -> Refactor garantizado"},
	{Key: "atomic-commits", Label: "Commits Atómicos", Description: "Límite estricto de 300 líneas por tarea"},
}

var namingOptions = []OptionItem{
	{Key: "standard", Label: "Estándar del Lenguaje", Description: "Convención canónica oficial"},
	{Key: "snake_case", Label: "snake_case", Description: "Archivos y elementos en snake_case"},
	{Key: "pascal_case", Label: "PascalCase", Description: "Clases y archivos en PascalCase"},
	{Key: "camel_case", Label: "camelCase", Description: "Notación camelCase"},
	{Key: "kebab_case", Label: "kebab-case", Description: "Archivos en kebab-case"},
}

var memoryOptions = []OptionItem{
	{Key: "sqlite-fts5", Label: "SQLite FTS5 Local", Description: "Búsqueda semántica y texto completo en .harness/memory.db"},
	{Key: "markdown", Label: "Archivos Markdown", Description: "Memoria plana en docs/decisions/"},
	{Key: "none", Label: "Sin Memoria Persistente", Description: "Solo reglas de agente sin indexación"},
}
