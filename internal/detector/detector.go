package detector

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type defaultDetector struct{}

// New crea una nueva instancia del detector por defecto.
func New() Detector {
	return &defaultDetector{}
}

func (d *defaultDetector) DetectDir(dirPath string) (*StackInfo, error) {
	return d.DetectFS(os.DirFS(dirPath))
}

func (d *defaultDetector) DetectFS(sysFS fs.FS) (*StackInfo, error) {
	files := make([]string, 0)

	// Explorar hasta un nivel razonable sin descender en directorios ignorados
	err := fs.WalkDir(sysFS, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == ".git" || name == "node_modules" || name == "bin" || name == "obj" || name == "vendor" || name == ".venv" {
				return fs.SkipDir
			}
			// Evitar descender más de 3 niveles
			if strings.Count(path, "/") > 3 {
				return fs.SkipDir
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 1. Evaluar C#
	if info := detectCSharp(sysFS, files); info != nil {
		return info, nil
	}

	// 2. Evaluar TypeScript / JavaScript
	if info := detectTypeScript(sysFS, files); info != nil {
		return info, nil
	}

	// 3. Evaluar Python
	if info := detectPython(sysFS, files); info != nil {
		return info, nil
	}

	// 4. Evaluar Go
	if info := detectGo(sysFS, files); info != nil {
		return info, nil
	}

	// 5. Sin coincidencias
	return &StackInfo{
		Language:   LangUnknown,
		TestRunner: RunnerUnknown,
		Confidence: 0.0,
		Evidence:   []string{},
	}, nil
}

func detectCSharp(sysFS fs.FS, files []string) *StackInfo {
	var csprojFiles []string
	var slnFiles []string

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".csproj" {
			csprojFiles = append(csprojFiles, file)
		} else if ext == ".sln" {
			slnFiles = append(slnFiles, file)
		}
	}

	if len(csprojFiles) == 0 && len(slnFiles) == 0 {
		return nil
	}

	info := &StackInfo{
		Language:   LangCSharp,
		TestRunner: RunnerUnknown,
		Confidence: 0.95,
		Evidence:   append(slnFiles, csprojFiles...),
	}

	targetFrameworkRegex := regexp.MustCompile(`(?i)<TargetFramework>([^<]+)</TargetFramework>`)

	for _, csproj := range csprojFiles {
		data, err := fs.ReadFile(sysFS, csproj)
		if err != nil {
			continue
		}
		content := string(data)
		lowerContent := strings.ToLower(content)

		if info.Framework == "" {
			matches := targetFrameworkRegex.FindStringSubmatch(content)
			if len(matches) > 1 {
				info.Framework = strings.TrimSpace(matches[1])
			}
		}

		if strings.Contains(lowerContent, "xunit") {
			info.TestRunner = RunnerXUnit
		} else if strings.Contains(lowerContent, "nunit") && info.TestRunner == RunnerUnknown {
			info.TestRunner = RunnerNUnit
		} else if strings.Contains(lowerContent, "mstest") && info.TestRunner == RunnerUnknown {
			info.TestRunner = RunnerMSTest
		}
	}

	return info
}

type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

func detectTypeScript(sysFS fs.FS, files []string) *StackInfo {
	hasPackageJSON := false
	hasTSConfig := false
	var evidence []string

	for _, file := range files {
		base := filepath.Base(file)
		if base == "package.json" {
			hasPackageJSON = true
			evidence = append(evidence, file)
		} else if base == "tsconfig.json" {
			hasTSConfig = true
			evidence = append(evidence, file)
		}
	}

	if !hasPackageJSON && !hasTSConfig {
		return nil
	}

	info := &StackInfo{
		Language:   LangTypeScript,
		TestRunner: RunnerUnknown,
		Confidence: 0.90,
		Evidence:   evidence,
	}

	data, err := fs.ReadFile(sysFS, "package.json")
	if err == nil {
		var pkg packageJSON
		if err := json.Unmarshal(data, &pkg); err == nil {
			allDeps := make(map[string]string)
			for k, v := range pkg.Dependencies {
				allDeps[strings.ToLower(k)] = v
			}
			for k, v := range pkg.DevDependencies {
				allDeps[strings.ToLower(k)] = v
			}

			// Test runner
			if _, ok := allDeps["vitest"]; ok {
				info.TestRunner = RunnerVitest
			} else if _, ok := allDeps["jest"]; ok {
				info.TestRunner = RunnerJest
			}

			// Framework
			if _, ok := allDeps["react"]; ok {
				info.Framework = "react"
			} else if _, ok := allDeps["@nestjs/core"]; ok {
				info.Framework = "nestjs"
			} else if _, ok := allDeps["next"]; ok {
				info.Framework = "nextjs"
			} else if _, ok := allDeps["express"]; ok {
				info.Framework = "express"
			}
		}
	}

	return info
}

func detectPython(sysFS fs.FS, files []string) *StackInfo {
	var evidence []string
	hasPy := false

	for _, file := range files {
		base := filepath.Base(file)
		if base == "pyproject.toml" || base == "requirements.txt" || base == "Pipfile" {
			evidence = append(evidence, file)
			hasPy = true
		}
	}

	if !hasPy {
		return nil
	}

	info := &StackInfo{
		Language:   LangPython,
		TestRunner: RunnerUnknown,
		Confidence: 0.90,
		Evidence:   evidence,
	}

	for _, ev := range evidence {
		data, err := fs.ReadFile(sysFS, ev)
		if err != nil {
			continue
		}
		content := strings.ToLower(string(data))

		if strings.Contains(content, "pytest") {
			info.TestRunner = RunnerPytest
		}

		if strings.Contains(content, "fastapi") {
			info.Framework = "fastapi"
		} else if strings.Contains(content, "django") && info.Framework == "" {
			info.Framework = "django"
		} else if strings.Contains(content, "flask") && info.Framework == "" {
			info.Framework = "flask"
		}
	}

	return info
}

func detectGo(sysFS fs.FS, files []string) *StackInfo {
	hasGoMod := false
	hasTest := false
	var evidence []string

	for _, file := range files {
		base := filepath.Base(file)
		if base == "go.mod" {
			hasGoMod = true
			evidence = append(evidence, file)
		} else if strings.HasSuffix(file, "_test.go") {
			hasTest = true
			evidence = append(evidence, file)
		}
	}

	if !hasGoMod {
		return nil
	}

	info := &StackInfo{
		Language:   LangGo,
		TestRunner: RunnerUnknown,
		Confidence: 0.95,
		Evidence:   evidence,
	}

	if hasTest {
		info.TestRunner = RunnerGoTest
	}

	data, err := fs.ReadFile(sysFS, "go.mod")
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		goVersionRegex := regexp.MustCompile(`^go\s+([0-9]+(\.[0-9]+)*)`)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			matches := goVersionRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				info.Framework = "go" + matches[1]
				break
			}
		}
	}

	return info
}
