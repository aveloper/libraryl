package home

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AddRoute(r *mux.Router, h *Handler) {
	r.HandleFunc("/", h.homePage).Methods(http.MethodGet)
}
