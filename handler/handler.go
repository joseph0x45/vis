package handler

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.renderDashboard)
}
