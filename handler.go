package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type handler struct {
	store   *Store
	baseURL string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/shorten":
		h.shorten(w, r)
	case r.Method == http.MethodGet && r.URL.Path != "/":
		h.redirect(w, r)
	default:
		http.NotFound(w, r)
	}
}

type shortenReq struct {
	URL string `json:"url"`
}

type shortenRes struct {
	Slug  string `json:"slug"`
	Short string `json:"short"`
}

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "url is required"})
		return
	}

	slug, err := h.store.Set(req.URL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate slug"})
		return
	}

	writeJSON(w, http.StatusOK, shortenRes{
		Slug:  slug,
		Short: h.baseURL + "/" + slug,
	})
}

func (h *handler) redirect(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := h.store.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
