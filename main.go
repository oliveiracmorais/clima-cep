package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/oliveiracmorais/cep-clima/handler"
)

//go:embed templates/*
var templateFiles embed.FS

func main() {
	cepHandler := handler.NewCEPHandler()

	// Carrega o template
	tmpl := template.Must(template.New("index.html").ParseFS(templateFiles, "templates/*"))

	http.HandleFunc("/cep/", cepHandler.GetClimaPorCEP)
	http.HandleFunc("/health", cepHandler.HealthCheck)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "index.html", nil)

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor iniciado na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
