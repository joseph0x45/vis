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

func makeChartData(readings []models.Reading) map[string]any {
	labels := make([]string, len(readings))
	values := make([]float64, len(readings))

	for i, reading := range readings {
		labels[i] = reading.DateString
		values[i] = reading.Kwh
	}
	chartData := map[string]any{
		"labels": labels,
		"values": values,
	}
	return chartData
}

func (h *Handler) renderDashboard(w http.ResponseWriter, r *http.Request) {
	latestReading, err := h.conn.GetLatestReading()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if latestReading == nil {
		latestReading = &models.Reading{
			DateString: "No Data yet",
		}
	}
	readings, err := h.conn.GetReadings(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chartData := makeChartData(readings)
	var contentBuffer bytes.Buffer
	pageData := models.PageData{
		ChartData:   chartData,
		LastUpdated: latestReading.DateString,
		CurrentKWh:  latestReading.Kwh,
		Timestamp:   latestReading.Timestamp,
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
		DateString: now.Format("January 2, 2006"),
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
		DateString: now.Format("January 2, 2006"),
	}
	if err := h.conn.InsertPurchase(purchase); err != nil {
		log.Println(errMsg + err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) renderHistory(w http.ResponseWriter, r *http.Request) {
	var err error
	var readings []models.Reading
	var purchases []models.Purchase
	if readings, err = h.conn.GetReadings(false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if purchases, err = h.conn.GetPurchases(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var contentBuffer bytes.Buffer
	if err = h.templates.ExecuteTemplate(&contentBuffer, "history", map[string]any{
		"Readings":  readings,
		"Purchases": purchases,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var responseBuffer bytes.Buffer
	pageData := models.PageData{
		Content: template.HTML(contentBuffer.String()),
	}
	if err = h.templates.ExecuteTemplate(&responseBuffer, "base", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBuffer.Bytes())
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.renderDashboard)
	r.Get("/history", h.renderHistory)
	r.Post("/readings", h.recordReading)
	r.Post("/purchases", h.logPurchase)
}
