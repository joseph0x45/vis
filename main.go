package main

import (
	"embed"
	"flag"
	"fmt"
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
	funcMap := template.FuncMap{
		"formatTime": func(timestamp int64) string {
			t := time.Unix(timestamp, 0)
			return t.Format("3:04 PM")
		},
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseFS(
		templatesFS,
		"templates/*.html",
	))
}

var version = "dev"

func main() {
	goutils.SetAppName("vis")
	versionFlag := flag.Bool("version", false, "Display the current version")
	generateServiceFileFlag := flag.Bool("generate-service-file", false, "Generate a service file")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Vis %s\n", version)
		return
	}
	if *generateServiceFileFlag {
		goutils.GenerateServiceFile("Vis, self hosted, personal power usage tracker")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
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
	log.Printf("Starting server on  http://0.0.0.0:%s", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
