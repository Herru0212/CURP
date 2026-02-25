package logic

import (
	"strings"
)

// GenerarCURP procesa los datos básicos
func GenerarCURP(nombre, paterno, materno, fecha, sexo, estado string) string {
	res := ""
	// 1. Inicial del primer apellido y primera vocal interna
	res += string(paterno[0])
	res += extraerVocalInterna(paterno)

	// 2. Inicial del segundo apellido
	res += string(materno[0])

	// 3. Inicial del primer nombre
	res += string(nombre[0])

	// Continuar con fecha (AAMMDD) y demás siglas
	return strings.ToUpper(res)
}
