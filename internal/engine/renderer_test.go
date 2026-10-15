package engine_test

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/FakuInsa/harness-up/assets"
	"github.com/FakuInsa/harness-up/internal/config"
	"github.com/FakuInsa/harness-up/internal/engine"
)

func TestRenderer_RenderString_FuncMap(t *testing.T) {
	r := engine.NewRenderer(fstest.MapFS{})

	tests := []struct {
		name     string
		tmpl     string
		data     any
		expected string
	}{
		{
			name:     "toLower and toUpper",
			tmpl:     "{{toLower .Text}} - {{toUpper .Text}}",
			data:     map[string]string{"Text": "HelloWorld"},
			expected: "helloworld - HELLOWORLD",
		},
		{
			name:     "camelCase from kebab",
			tmpl:     "{{camelCase .}}",
			data:     "order-processing-service",
			expected: "orderProcessingService",
		},
		{
			name:     "camelCase from snake",
			tmpl:     "{{camelCase .}}",
			data:     "user_profile_dto",
			expected: "userProfileDto",
		},
		{
			name:     "pascalCase from kebab",
			tmpl:     "{{pascalCase .}}",
			data:     "order-processing-service",
			expected: "OrderProcessingService",
		},
		{
			name:     "pascalCase from snake",
			tmpl:     "{{pascalCase .}}",
			data:     "user_profile_dto",
			expected: "UserProfileDto",
		},
		{
			name:     "snakeCase from camel and pascal",
			tmpl:     "{{snakeCase .Camel}} - {{snakeCase .Pascal}}",
			data:     map[string]string{"Camel": "helloWorldCase", "Pascal": "OrderProcessorService"},
			expected: "hello_world_case - order_processor_service",
		},
		{
			name:     "kebabCase from camel and pascal",
			tmpl:     "{{kebabCase .Camel}} - {{kebabCase .Pascal}}",
			data:     map[string]string{"Camel": "helloWorldCase", "Pascal": "OrderProcessorService"},
			expected: "hello-world-case - order-processor-service",
		},
		{
			name:     "contains string",
			tmpl:     "{{if contains .List \"antigravity\"}}found{{else}}missing{{end}}",
			data:     map[string][]string{"List": {"cursor", "antigravity", "claude"}},
			expected: "found",
		},
		{
			name:     "contains missing",
			tmpl:     "{{if contains .List \"windsurf\"}}found{{else}}missing{{end}}",
			data:     map[string][]string{"List": {"cursor", "antigravity"}},
			expected: "missing",
		},
		{
			name:     "join strings",
			tmpl:     "{{join .List \", \"}}",
			data:     map[string][]string{"List": {"go", "csharp", "typescript"}},
			expected: "go, csharp, typescript",
		},
		{
			name:     "empty casing and symbols",
			tmpl:     "A:{{camelCase .}} B:{{pascalCase .}} C:{{snakeCase .}} D:{{kebabCase .}}",
			data:     "---",
			expected: "A: B: C: D:",
		},
		{
			name:     "contains in string substring and slice of any",
			tmpl:     "{{if contains .Str \"world\"}}hasStr{{end}} {{if contains .AnySlice \"test\"}}hasAny{{end}} {{if contains 123 \"x\"}}hasNone{{else}}no{{end}}",
			data:     map[string]any{"Str": "hello world", "AnySlice": []any{"foo", "test", 42}},
			expected: "hasStr hasAny no",
		},
		{
			name:     "join slice of any and scalar fallback",
			tmpl:     "{{join .AnyList \"-\"}} / {{join .Scalar \",\"}}",
			data:     map[string]any{"AnyList": []any{"alpha", 100, true}, "Scalar": "single"},
			expected: "alpha-100-true / single",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := r.RenderString(tc.tmpl, tc.data)
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}
			if out != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, out)
			}
		})
	}
}

func TestRenderer_RenderString_Conditionals(t *testing.T) {
	r := engine.NewRenderer(fstest.MapFS{})

	tmpl := `{{if eq .Pattern "vertical-slice"}}Features First{{else}}Layered{{end}}`

	out1, err := r.RenderString(tmpl, map[string]string{"Pattern": "vertical-slice"})
	if err != nil || out1 != "Features First" {
		t.Fatalf("expected 'Features First', got %q, err: %v", out1, err)
	}

	out2, err := r.RenderString(tmpl, map[string]string{"Pattern": "clean"})
	if err != nil || out2 != "Layered" {
		t.Fatalf("expected 'Layered', got %q, err: %v", out2, err)
	}
}

func TestRenderer_RenderString_InvalidSyntax(t *testing.T) {
	r := engine.NewRenderer(fstest.MapFS{})

	_, err := r.RenderString("Hello {{if .Name}} broken template", map[string]string{"Name": "Test"})
	if err == nil {
		t.Fatal("expected syntax error on unclosed template block, got nil")
	}
}

func TestRenderer_RenderAsset_MapFS(t *testing.T) {
	mockFS := fstest.MapFS{
		"rules/test.md.tmpl": &fstest.MapFile{
			Data: []byte("Project: {{.Project.Name}} ({{kebabCase .Project.Name}})"),
		},
	}

	r := engine.NewRenderer(mockFS)
	cfg := config.Default("ECommerce Core")

	out, err := r.RenderAsset("rules/test.md.tmpl", cfg)
	if err != nil {
		t.Fatalf("unexpected error rendering asset: %v", err)
	}

	expected := "Project: ECommerce Core (e-commerce-core)"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestRenderer_MissingAsset(t *testing.T) {
	r := engine.NewRenderer(fstest.MapFS{})

	_, err := r.RenderAsset("rules/non_existent.md.tmpl", nil)
	if err == nil {
		t.Fatal("expected error when asset is missing, got nil")
	}
}

func TestRenderer_ListAssets(t *testing.T) {
	mockFS := fstest.MapFS{
		"rules/arch.md.tmpl":       &fstest.MapFile{Data: []byte("content")},
		"rules/guard.md.tmpl":      &fstest.MapFile{Data: []byte("content")},
		"skills/workflow/SKILL.md": &fstest.MapFile{Data: []byte("content")},
		"other/random.txt":         &fstest.MapFile{Data: []byte("content")},
	}

	r := engine.NewRenderer(mockFS)

	rules, err := r.ListAssets("rules")
	if err != nil {
		t.Fatalf("unexpected error listing rules: %v", err)
	}

	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d: %v", len(rules), rules)
	}

	skills, err := r.ListAssets("skills")
	if err != nil {
		t.Fatalf("unexpected error listing skills: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d: %v", len(skills), skills)
	}
}

func TestRenderer_RenderAsset_EmbeddedFS(t *testing.T) {
	r := engine.NewRenderer(assets.FS)
	cfg := config.Default("TestApp")

	// Render architecture template
	archOutput, err := r.RenderAsset("rules/architecture.md.tmpl", cfg)
	if err != nil {
		t.Fatalf("failed to render embedded architecture.md.tmpl: %v", err)
	}
	if !strings.Contains(archOutput, "vertical-slice") && !strings.Contains(archOutput, "Vertical Slice") {
		t.Errorf("expected architecture output to contain vertical-slice directives, got: %s", archOutput)
	}

	// Render guardrails template
	guardOutput, err := r.RenderAsset("rules/guardrails.md.tmpl", cfg)
	if err != nil {
		t.Fatalf("failed to render embedded guardrails.md.tmpl: %v", err)
	}
	if !strings.Contains(guardOutput, "300") {
		t.Errorf("expected guardrails output to mention max diff 300, got: %s", guardOutput)
	}
}

type failData struct{}

func (failData) CauseError() (string, error) {
	return "", fmt.Errorf("forced execution error")
}

func TestRenderer_RenderString_ExecuteError(t *testing.T) {
	r := engine.NewRenderer(fstest.MapFS{})
	_, err := r.RenderString("{{.CauseError}}", failData{})
	if err == nil {
		t.Fatal("expected execution error when calling error-returning method, got nil")
	}
}
