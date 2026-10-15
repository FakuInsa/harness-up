package main

import (
	"fmt"
	"os"

	"github.com/FakuInsa/harness-up/internal/wizard"
)

func main() {
	fmt.Println("Iniciando HarnessUp Wizard...")

	cfg, err := wizard.Run(".")
	if err != nil {
		fmt.Printf("\n❌ Wizard cancelado o con error: %v\n", err)
		os.Exit(1)
	}

	yamlData, err := cfg.ToYAML()
	if err != nil {
		fmt.Printf("\n❌ Error serializando configuración: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Configuración generada exitosamente:\n\n%s\n", string(yamlData))
}
