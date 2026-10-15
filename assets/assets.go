package assets

import "embed"

// FS contiene todos los archivos estáticos de reglas, skills y templates embebidos en el binario.
//
//go:embed rules/* skills/* templates/*
var FS embed.FS
