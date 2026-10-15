package emitter

import (
	"os"
	"path/filepath"
	"strings"
)

// MergeContent fusiona el contenido gestionado dentro de un archivo existente respetando delimitadores y zonas de usuario.
func MergeContent(existingContent, newManagedContent string) string {
	return MergeWithHeader(existingContent, "", newManagedContent)
}

// MergeWithHeader fusiona el contenido gestionado permitiendo anteponer un encabezado fijo (ej. frontmatter en Cursor).
func MergeWithHeader(existingContent, header, newManagedContent string) string {
	var sb strings.Builder

	startIndex := strings.Index(existingContent, ManagedStart)
	endIndex := strings.Index(existingContent, ManagedEnd)

	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		// Archivo con marcas existentes
		if header != "" {
			sb.WriteString(strings.TrimSpace(header))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(existingContent[:startIndex])
		}

		sb.WriteString(ManagedStart)
		sb.WriteString("\n<!-- No edites este bloque manualmente. Es gestionado por HarnessUp. -->\n")
		sb.WriteString(strings.TrimSpace(newManagedContent))
		sb.WriteString("\n")
		sb.WriteString(ManagedEnd)

		postamble := existingContent[endIndex+len(ManagedEnd):]
		if !strings.Contains(postamble, CustomUserRulesTag) {
			sb.WriteString("\n\n")
			sb.WriteString(CustomUserRulesTag)
			sb.WriteString("\n<!-- Agrega aquí tus reglas personalizadas que no deban ser sobreescritas por HarnessUp. -->\n")
		}
		sb.WriteString(postamble)
		return sb.String()
	}

	// Archivo sin marcas previas
	if header != "" {
		sb.WriteString(strings.TrimSpace(header))
		sb.WriteString("\n\n")
	}

	sb.WriteString(ManagedStart)
	sb.WriteString("\n<!-- No edites este bloque manualmente. Es gestionado por HarnessUp. -->\n")
	sb.WriteString(strings.TrimSpace(newManagedContent))
	sb.WriteString("\n")
	sb.WriteString(ManagedEnd)
	sb.WriteString("\n\n")
	sb.WriteString(CustomUserRulesTag)
	sb.WriteString("\n")

	trimmedExisting := strings.TrimSpace(existingContent)
	if trimmedExisting != "" {
		sb.WriteString(trimmedExisting)
		sb.WriteString("\n")
	} else {
		sb.WriteString("<!-- Agrega aquí tus reglas personalizadas que no deban ser sobreescritas por HarnessUp. -->\n")
	}

	return sb.String()
}

// WriteMergedFile escribe un archivo fusionando su contenido con cualquier versión previa existente en disco.
func WriteMergedFile(filePath, newManagedContent string, perm os.FileMode) error {
	return WriteMergedFileWithHeader(filePath, "", newManagedContent, perm)
}

// WriteMergedFileWithHeader escribe un archivo con cabecera y contenido fusionado de forma idempotente.
func WriteMergedFileWithHeader(filePath, header, newManagedContent string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var existingContent string
	data, err := os.ReadFile(filePath)
	if err == nil {
		existingContent = string(data)
	} else if !os.IsNotExist(err) {
		return err
	}

	merged := MergeWithHeader(existingContent, header, newManagedContent)
	return os.WriteFile(filePath, []byte(merged), perm)
}
