package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/sad"
	"github.com/joseph0x45/vis/db"
	"github.com/joseph0x45/vis/handler"
)

//go:embed output.css
var stylesCSS template.CSS

//go:embed chart.js
var chartJS template.JS

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(
		templatesFS,
		"templates/*.html",
	))
}

func main() {
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}
	goutils.SetAppName("vis")
	dbPath := goutils.Setup()
	conn := db.Connect(sad.DBConnectionOptions{
		EnableForeignKeys: true,
		DatabasePath:      dbPath,
	})
	r := chi.NewRouter()
	handler := handler.NewHandler(templates, conn)
	server := http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	r.Get("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(stylesCSS))
	})
	r.Get("/chart.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(chartJS))
	})
	handler.RegisterRoutes(r)
	log.Printf("Starting server on  http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
