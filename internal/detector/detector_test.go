package detector_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/FakuInsa/harness-up/internal/detector"
)

func TestDetect_CSharp_XUnit(t *testing.T) {
	mockFS := fstest.MapFS{
		"Solution.sln": &fstest.MapFile{Data: []byte("Microsoft Visual Studio Solution File...")},
		"src/Api/Api.csproj": &fstest.MapFile{
			Data: []byte(`
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net9.0</TargetFramework>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="xunit" Version="2.9.0" />
    <PackageReference Include="xunit.runner.visualstudio" Version="2.8.2" />
  </ItemGroup>
</Project>`),
		},
	}

	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangCSharp {
		t.Errorf("expected %s, got %s", detector.LangCSharp, info.Language)
	}
	if info.TestRunner != detector.RunnerXUnit {
		t.Errorf("expected %s, got %s", detector.RunnerXUnit, info.TestRunner)
	}
	if info.Framework != "net9.0" {
		t.Errorf("expected framework net9.0, got %s", info.Framework)
	}
	if len(info.Evidence) == 0 {
		t.Errorf("expected evidence files, got none")
	}
}

func TestDetect_TypeScript_Vitest(t *testing.T) {
	mockFS := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{
  "name": "my-web-app",
  "dependencies": {
    "react": "^19.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "vitest": "^2.0.0"
  }
}`),
		},
		"tsconfig.json": &fstest.MapFile{Data: []byte("{}")},
	}

	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangTypeScript {
		t.Errorf("expected %s, got %s", detector.LangTypeScript, info.Language)
	}
	if info.TestRunner != detector.RunnerVitest {
		t.Errorf("expected %s, got %s", detector.RunnerVitest, info.TestRunner)
	}
	if info.Framework != "react" {
		t.Errorf("expected framework react, got %s", info.Framework)
	}
}

func TestDetect_Python_Pytest(t *testing.T) {
	mockFS := fstest.MapFS{
		"pyproject.toml": &fstest.MapFile{
			Data: []byte(`
[project]
name = "api-service"
dependencies = [
    "fastapi>=0.110.0"
]

[tool.pytest.ini_options]
minversion = "7.0"
`),
		},
	}

	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangPython {
		t.Errorf("expected %s, got %s", detector.LangPython, info.Language)
	}
	if info.TestRunner != detector.RunnerPytest {
		t.Errorf("expected %s, got %s", detector.RunnerPytest, info.TestRunner)
	}
	if info.Framework != "fastapi" {
		t.Errorf("expected framework fastapi, got %s", info.Framework)
	}
}

func TestDetect_Go_Standard(t *testing.T) {
	mockFS := fstest.MapFS{
		"go.mod": &fstest.MapFile{
			Data: []byte("module github.com/foo/bar\n\ngo 1.24.0\n"),
		},
		"main_test.go": &fstest.MapFile{
			Data: []byte("package main\nimport \"testing\"\nfunc TestFoo(t *testing.T){}"),
		},
	}

	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangGo {
		t.Errorf("expected %s, got %s", detector.LangGo, info.Language)
	}
	if info.TestRunner != detector.RunnerGoTest {
		t.Errorf("expected %s, got %s", detector.RunnerGoTest, info.TestRunner)
	}
	if info.Framework != "go1.24.0" {
		t.Errorf("expected framework go1.24.0, got %s", info.Framework)
	}
}

func TestDetect_EmptyOrUnknown(t *testing.T) {
	mockFS := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# Empty Project")},
	}

	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangUnknown {
		t.Errorf("expected %s, got %s", detector.LangUnknown, info.Language)
	}
	if info.TestRunner != detector.RunnerUnknown {
		t.Errorf("expected %s, got %s", detector.RunnerUnknown, info.TestRunner)
	}
	if info.Confidence != 0.0 {
		t.Errorf("expected confidence 0.0, got %f", info.Confidence)
	}
}

func TestDetect_CSharp_NUnit_MSTest(t *testing.T) {
	mockNUnit := fstest.MapFS{
		"App.csproj": &fstest.MapFile{
			Data: []byte(`<Project><ItemGroup><PackageReference Include="nunit" Version="3.13" /></ItemGroup></Project>`),
		},
	}
	d := detector.New()
	info, err := d.DetectFS(mockNUnit)
	if err != nil || info.TestRunner != detector.RunnerNUnit {
		t.Fatalf("expected nunit, got %v (err: %v)", info.TestRunner, err)
	}

	mockMSTest := fstest.MapFS{
		"App.csproj": &fstest.MapFile{
			Data: []byte(`<Project><ItemGroup><PackageReference Include="MSTest.TestFramework" Version="3.0" /></ItemGroup></Project>`),
		},
	}
	info, err = d.DetectFS(mockMSTest)
	if err != nil || info.TestRunner != detector.RunnerMSTest {
		t.Fatalf("expected mstest, got %v (err: %v)", info.TestRunner, err)
	}
}

func TestDetect_TypeScript_Jest_NestJS(t *testing.T) {
	mockFS := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{
  "dependencies": { "@nestjs/core": "^10.0.0" },
  "devDependencies": { "jest": "^29.0.0" }
}`),
		},
	}
	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil || info.TestRunner != detector.RunnerJest || info.Framework != "nestjs" {
		t.Fatalf("expected jest and nestjs, got runner: %v, fw: %v (err: %v)", info.TestRunner, info.Framework, err)
	}
}

func TestDetect_Python_Django_Flask(t *testing.T) {
	mockFS := fstest.MapFS{
		"requirements.txt": &fstest.MapFile{
			Data: []byte("django>=5.0\npytest\n"),
		},
	}
	d := detector.New()
	info, err := d.DetectFS(mockFS)
	if err != nil || info.Framework != "django" || info.TestRunner != detector.RunnerPytest {
		t.Fatalf("expected django and pytest, got fw: %v, runner: %v (err: %v)", info.Framework, info.TestRunner, err)
	}

	mockFlask := fstest.MapFS{
		"Pipfile": &fstest.MapFile{
			Data: []byte("flask = \"*\"\n"),
		},
	}
	info, err = d.DetectFS(mockFlask)
	if err != nil || info.Framework != "flask" {
		t.Fatalf("expected flask, got fw: %v (err: %v)", info.Framework, err)
	}
}

func TestDetect_DetectDir(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module temp\n\ngo 1.23.0\n"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	d := detector.New()
	info, err := d.DetectDir(tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Language != detector.LangGo {
		t.Errorf("expected %s, got %s", detector.LangGo, info.Language)
	}
}
