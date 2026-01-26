package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/vis/db"
	"github.com/joseph0x45/vis/models"
)

type Handler struct {
	templates *template.Template
	conn      *db.Conn
	version   string
}

func NewHandler(
	templates *template.Template,
	conn *db.Conn,
	version string,
) *Handler {
	return &Handler{
		templates: templates,
		conn:      conn,
		version:   version,
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

func (h *Handler) renderLoginPage(w http.ResponseWriter, r *http.Request) {
	var responseBuffer bytes.Buffer
	if err := h.templates.ExecuteTemplate(&responseBuffer, "login", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBuffer.Bytes())
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.templates.ExecuteTemplate(w, "login", map[string]any{
			"Error": "Something went wrong! Try again",
		})
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.conn.GetUserBy("username", username)
	if err != nil {
		h.templates.ExecuteTemplate(w, "login", map[string]any{
			"Error": "Something went wrong! Try again",
		})
		return
	}
	if user == nil {
		h.templates.ExecuteTemplate(w, "login", map[string]any{
			"Error": fmt.Sprintf("User %s not found", username),
		})
		return
	}
	if !goutils.HashMatchesPassword(user.Password, password) {
		h.templates.ExecuteTemplate(w, "login", map[string]any{
			"Error": "Invalid credentials",
		})
		return
	}
	cookie := &http.Cookie{
		Name:     "user",
		Value:    user.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.version != "dev",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(100 * 365 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		if err != nil {
			log.Println("Failed to get cookie:", err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := h.conn.GetUserBy("id", cookie.Value)
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/login", h.renderLoginPage)
	r.Post("/login", h.login)
	r.Group(func(r chi.Router) {
		r.Use(h.authRequired)
		r.Get("/", h.renderDashboard)
		r.Get("/history", h.renderHistory)
		r.Post("/readings", h.recordReading)
		r.Post("/purchases", h.logPurchase)
	})
}
