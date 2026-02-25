package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ANALIZADOR LÉXICO: Funciones para extraer caracteres específicos
func esVocal(char byte) bool {
	return strings.ContainsAny(string(char), "AEIOU")
}

func primeraVocalInterna(s string) string {
	for i := 1; i < len(s); i++ {
		if esVocal(s[i]) {
			return string(s[i])
		}
	}
	return "X"
}

func primeraConsonanteInterna(s string) string {
	for i := 1; i < len(s); i++ {
		// Si no es vocal y es una letra de la A-Z, es consonante
		if !esVocal(s[i]) && s[i] >= 'A' && s[i] <= 'Z' {
			return string(s[i])
		}
	}
	return "X"
}

func main() {
	// Servir la interfaz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	// Procesar la CURP (Analizador Sintáctico)
	http.HandleFunc("/generar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Captura de datos
			nombre := strings.ToUpper(r.FormValue("nombre"))
			paterno := strings.ToUpper(r.FormValue("paterno"))
			materno := strings.ToUpper(r.FormValue("materno"))
			fecha := strings.ReplaceAll(r.FormValue("fecha"), "-", "") // Formato YYYYMMDD
			genero := r.FormValue("genero")
			estado := r.FormValue("estado")

			// CONSTRUCCIÓN SINTÁCTICA DE LA CURP (18 Caracteres)
			var res strings.Builder

			// 1. Inicial y Vocal del Paterno
			res.WriteByte(paterno[0])
			res.WriteString(primeraVocalInterna(paterno))

			// 2. Inicial Materno
			if len(materno) > 0 {
				res.WriteByte(materno[0])
			} else {
				res.WriteString("X")
			}

			// 3. Inicial Nombre
			res.WriteByte(nombre[0])

			// 4. Fecha (AAMMDD)
			res.WriteString(fecha[2:])

			// 5. Género (H/M)
			res.WriteString(genero)

			// 6. Estado (TS / TL)
			res.WriteString(estado)

			// 7. Consonantes Internas (Paterno, Materno, Nombre)
			res.WriteString(primeraConsonanteInterna(paterno))
			if len(materno) > 0 {
				res.WriteString(primeraConsonanteInterna(materno))
			} else {
				res.WriteString("X")
			}
			res.WriteString(primeraConsonanteInterna(nombre))

			// 8. Homoclave y Verificador (Simulados para el avance)
			res.WriteString("01")

			// Responder solo con el texto de la CURP
			fmt.Fprint(w, res.String())
		}
	})

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
