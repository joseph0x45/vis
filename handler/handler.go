package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joseph0x45/vis/db"
	"github.com/joseph0x45/vis/models"
)

type Handler struct {
	templates *template.Template
	conn      *db.Conn
}

func NewHandler(
	templates *template.Template,
	conn *db.Conn,
) *Handler {
	return &Handler{
		templates: templates,
		conn:      conn,
	}
}

func (h *Handler) renderDashboard(w http.ResponseWriter, r *http.Request) {
	var contentBuffer bytes.Buffer
	pageData := models.PageData{
		LastUpdated: "Yesterday",
		CurrentKWh:  143,
	}
	if err := h.templates.ExecuteTemplate(&contentBuffer, "dashboard", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageData.Content = template.HTML(contentBuffer.String())
	var responseBuffer bytes.Buffer
	if err := h.templates.ExecuteTemplate(&responseBuffer, "base", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBuffer.Bytes())
}

func (h *Handler) recordReading(w http.ResponseWriter, r *http.Request) {
	const errMsg = "Error while recording reading: "
	if err := r.ParseForm(); err != nil {
		log.Println(errMsg+"Failed to parse form:", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	kwh, err := strconv.ParseFloat(r.FormValue("kwh"), 64)
	if err != nil {
		log.Println(errMsg+"Failed to convert kwh to float:", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	now := time.Now()
	reading := &models.Reading{
		ID:         uuid.NewString(),
		Timestamp:  now.Unix(),
		Kwh:        kwh,
		DateString: now.Format("02-01-2006"),
	}
	if err := h.conn.InsertReading(reading); err != nil {
		log.Println(errMsg + err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logPurchase(w http.ResponseWriter, r *http.Request) {
	const errMsg = "Error while logging purchase: "
	if err := r.ParseForm(); err != nil {
		log.Println(errMsg+"Failed to parse form:", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	kwh, err := strconv.ParseFloat(r.FormValue("kwh"), 64)
	if err != nil {
		log.Println(errMsg+"Failed to convert kwh to float:", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	cost, err := strconv.Atoi(r.FormValue("cost"))
	if err != nil {
		log.Println(errMsg+"Failed to convert cost to int:", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	now := time.Now()
	purchase := &models.Purchase{
		ID:         uuid.NewString(),
		Timestamp:  now.Unix(),
		Kwh:        kwh,
		Cost:       cost,
		DateString: now.Format("02-01-2006"),
	}
	if err := h.conn.InsertPurchase(purchase); err != nil {
		log.Println(errMsg + err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.renderDashboard)
	r.Post("/readings", h.recordReading)
	r.Post("/purchases", h.logPurchase)
}
