package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Persona struct {
	Nombre          string
	Paterno         string
	Materno         string
	FechaNacimiento string
	Genero          string
	Estado          string
	CURP            string
	Error           string
}

// Auxiliar: Busca primera vocal interna (después de la primera letra)
func primeraVocalInterna(s string) uint8 {
	vocales := "AEIOU"
	for i := 1; i < len(s); i++ {
		if strings.ContainsAny(string(s[i]), vocales) {
			return s[i]
		}
	}
	return 'X'
}

// Auxiliar: Busca primera consonante interna (después de la primera letra)
func primeraConsonanteInterna(s string) uint8 {
	consonantes := "BCDFGHJKLMNPQRSTVWXYZ"
	for i := 1; i < len(s); i++ {
		if strings.ContainsAny(string(s[i]), consonantes) {
			return s[i]
		}
	}
	return 'X'
}

// Regla especial para MARÍA y JOSÉ
func filtrarNombre(nombre string) string {
	partes := strings.Fields(nombre)
	if len(partes) > 1 {
		primero := partes[0]
		if primero == "MARIA" || primero == "JOSE" || primero == "MA" || primero == "MA." || primero == "J" || primero == "J." {
			return partes[1]
		}
	}
	return partes[0]
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Generador CURP Oficial</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #2c3e50; color: white; display: flex; justify-content: center; padding: 20px; }
        .container { background: white; color: #333; padding: 30px; border-radius: 15px; width: 450px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); }
        h2 { text-align: center; color: #1a5276; border-bottom: 2px solid #1a5276; padding-bottom: 10px; }
        input, select { width: 100%; padding: 12px; margin: 10px 0; border: 1px solid #ccc; border-radius: 8px; box-sizing: border-box; }
        button { width: 100%; padding: 15px; background: #1a5276; color: white; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; font-size: 1rem; }
        .result { margin-top: 20px; background: #d4efdf; padding: 15px; border-radius: 8px; border: 2px solid #27ae60; text-align: center; }
        .error { background: #fadbd8; color: #a93226; padding: 10px; border-radius: 8px; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Generador de CURP</h2>
        <form method="POST">
            <input type="text" name="nombre" placeholder="Nombre(s)" required>
            <input type="text" name="paterno" placeholder="Primer Apellido" required>
            <input type="text" name="materno" placeholder="Segundo Apellido (o X si no tiene)" required>
            <label>Fecha de Nacimiento:</label>
            <input type="date" name="fecha" required>
            <select name="genero">
                <option value="H">Hombre</option>
                <option value="M">Mujer</option>
                <option value="X">No Binario</option>
            </select>
            <select name="estado">
                <option value="TS">Tamaulipas</option>
                <option value="TL">Tlaxcala</option>
                <option value="NE">Nacido en el Extranjero</option>
            </select>
            <button type="submit">Generar CURP Completa</button>
        </form>

        {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
        {{if .CURP}}
            <div class="result">
                <strong>CURP Oficial Generada:</strong><br>
                <span style="font-size: 1.8rem; letter-spacing: 2px; font-family: monospace;">{{.CURP}}</span>
            </div>
        {{end}}
    </div>
</body>
</html>
`

func main() {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 1. Captura de Tokens
			nombreRaw := strings.ToUpper(r.FormValue("nombre"))
			paterno := strings.ToUpper(r.FormValue("paterno"))
			materno := strings.ToUpper(r.FormValue("materno"))
			fechaStr := r.FormValue("fecha")
			genero := r.FormValue("genero")
			estado := r.FormValue("estado")

			// 2. Validación de 150 años (Analizador Sintáctico)
			fechaNac, _ := time.Parse("2006-01-02", fechaStr)
			if time.Since(fechaNac).Hours() > 150*365*24 || fechaNac.After(time.Now()) {
				tmpl.Execute(w, Persona{Error: "Error: La fecha excede el límite de 150 años o es futura."})
				return
			}

			// 3. Aplicación de reglas de negocio
			nombreFiltrado := filtrarNombre(nombreRaw)

			// PARTE 1: Iniciales y Vocales
			p1 := fmt.Sprintf("%c%c%c%c",
				paterno[0],
				primeraVocalInterna(paterno),
				materno[0],
				nombreFiltrado[0])

			// PARTE 2: Fecha (AAMMDD)
			f := strings.ReplaceAll(fechaStr, "-", "")
			p2 := f[2:8]

			// PARTE 3: Sexo y Estado
			p3 := genero + estado

			// PARTE 4: Consonantes Internas
			p4 := fmt.Sprintf("%c%c%c",
				primeraConsonanteInterna(paterno),
				primeraConsonanteInterna(materno),
				primeraConsonanteInterna(nombreFiltrado))

			// PARTE 5: Homoclave y Verificador (Dígitos finales)
			p5 := "A1"

			curpFinal := p1 + p2 + p3 + p4 + p5
			tmpl.Execute(w, Persona{CURP: strings.ToUpper(curpFinal)})
			return
		}
		tmpl.Execute(w, nil)
	})

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
