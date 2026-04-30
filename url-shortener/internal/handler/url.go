package handler

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"url-shortener/internal/service"
)

type URLHandler struct {
	service   *service.URLService
	templates *template.Template
}

type HomePageData struct {
	Items []service.ShortURL
}

func NewURLHandler(service *service.URLService) *URLHandler {
	templates := template.Must(template.ParseFiles("web/templates/index.html"))

	return &URLHandler{
		service:   service,
		templates: templates,
	}
}

func (h *URLHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	items, err := h.service.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load urls"})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "index.html", HomePageData{Items: items}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to render page"})
	}
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req service.CreateShortURLRequest
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form body"})
			return
		}

		req.OriginalURL = r.FormValue("original_url")
	}

	shortURL, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid url"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create short url"})
		return
	}

	writeJSON(w, http.StatusCreated, shortURL)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing short code"})
		return
	}

	shortURL, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "short code not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to resolve short url"})
		return
	}

	http.Redirect(w, r, shortURL.OriginalURL, http.StatusFound)
}

func (h *URLHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/api/urls/")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing short code"})
		return
	}

	shortURL, err := h.service.GetByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "short code not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch short url"})
		return
	}

	writeJSON(w, http.StatusOK, shortURL)
}

func (h *URLHandler) ListShortURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	items, err := h.service.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list short urls"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *URLHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
