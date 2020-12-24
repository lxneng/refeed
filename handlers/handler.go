package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lxneng/refeed/models"
)

type Handler struct {
	Router *mux.Router
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

func (h *Handler) FeedHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	atom, err := models.GetAtomFeed(vars["slug"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error! %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, atom)
}
