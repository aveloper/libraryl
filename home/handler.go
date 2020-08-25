package home

import (
	"html/template"
	"net/http"
)

var homeTmpl = template.Must(template.ParseFiles("./views/home.html"))

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": "What brings you here homo-sepien?",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := homeTmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
