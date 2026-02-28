package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Persona struct {
	Nombre, Paterno, Materno, FechaNacimiento, Genero, Estado, CURP, Error string
}

func primeraVocalInterna(s string) uint8 {
	vocales := "AEIOU"
	for i := 1; i < len(s); i++ {
		if strings.ContainsAny(string(s[i]), vocales) {
			return s[i]
		}
	}
	return 'X'
}

func primeraConsonanteInterna(s string) uint8 {
	consonantes := "BCDFGHJKLMNPQRSTVWXYZ"
	for i := 1; i < len(s); i++ {
		if strings.ContainsAny(string(s[i]), consonantes) {
			return s[i]
		}
	}
	return 'X'
}

func filtrarNombre(nombre string) string {
	partes := strings.Fields(nombre)
	if len(partes) > 1 {
		p := partes[0]
		if p == "MARIA" || p == "JOSE" || p == "MA" || p == "MA." || p == "J" || p == "J." {
			return partes[1]
		}
	}
	return partes[0]
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Generador CURP Oficial</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #2c3e50; color: white; display: flex; justify-content: center; padding: 20px; }
        .container { background: white; color: #333; padding: 30px; border-radius: 15px; width: 450px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); }
        input, select { width: 100%; padding: 12px; margin: 10px 0; border: 1px solid #ccc; border-radius: 8px; box-sizing: border-box; }
        button { width: 100%; padding: 15px; background: #1a5276; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: bold; }
        .result { margin-top: 20px; background: #d4efdf; padding: 15px; border-radius: 8px; border: 2px solid #27ae60; text-align: center; }
        .error { background: #fadbd8; color: #a93226; padding: 10px; border-radius: 8px; margin-top: 10px; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h2 style="text-align:center; color:#1a5276">Generador de CURP</h2>
        <form method="POST">
            <input type="text" name="nombre" placeholder="Nombre(s)" required>
            <input type="text" name="paterno" placeholder="Primer Apellido" required>
            <input type="text" name="materno" placeholder="Segundo Apellido" required>
            <input type="date" name="fecha" required>
            <select name="genero"><option value="H">Hombre</option><option value="M">Mujer</option><option value="X">No Binario</option></select>
            <select name="estado"><option value="TS">Tamaulipas</option><option value="TL">Tlaxcala</option></select>
            <button type="submit">Generar CURP Completa</button>
        </form>
        {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
        {{if .CURP}}<div class="result"><strong>CURP Generada:</strong><br><span style="font-size:1.6rem; letter-spacing:2px; font-family:monospace;">{{.CURP}}</span></div>{{end}}
    </div>
</body>
</html>`

func main() {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	re := regexp.MustCompile(`^[a-zA-ZñÑ\s]+$`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			n, p, m := r.FormValue("nombre"), r.FormValue("paterno"), r.FormValue("materno")
			fStr := r.FormValue("fecha")

			// 1. Validar Letras y Ñ
			if !re.MatchString(n) || !re.MatchString(p) || !re.MatchString(m) {
				tmpl.Execute(w, Persona{Error: "Error: No se permiten números ni símbolos. Solo letras y Ñ."})
				return
			}
			// 2. Validar Longitud (Mínimo 3 letras)
			if len(p) < 3 || len(m) < 3 {
				tmpl.Execute(w, Persona{Error: "Error: Los apellidos deben tener al menos 3 letras."})
				return
			}
			// 3. Validar 150 años
			fNac, _ := time.Parse("2006-01-02", fStr)
			if time.Since(fNac).Hours() > 150*365*24 || fNac.After(time.Now()) {
				tmpl.Execute(w, Persona{Error: "Error: Fecha fuera de rango (Máximo 150 años)."})
				return
			}

			// Proceso de generación
			nomF := filtrarNombre(strings.ToUpper(n))
			pat, mat := strings.ToUpper(p), strings.ToUpper(m)
			f := strings.ReplaceAll(fStr, "-", "")

			res := fmt.Sprintf("%c%c%c%c%s%s%s%c%c%c01",
				pat[0], primeraVocalInterna(pat), mat[0], nomF[0],
				f[2:8], r.FormValue("genero"), r.FormValue("estado"),
				primeraConsonanteInterna(pat), primeraConsonanteInterna(mat), primeraConsonanteInterna(nomF))

			tmpl.Execute(w, Persona{CURP: res})
			return
		}
		tmpl.Execute(w, nil)
	})
	fmt.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
