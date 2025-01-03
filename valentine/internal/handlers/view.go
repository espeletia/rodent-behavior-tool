package handlers

import (
	"net/http"
	"valentine/view"
)

type ViewHandler struct {
}

func NewViewHandler() *ViewHandler {
	return &ViewHandler{}
}

func (vh *ViewHandler) Render(w http.ResponseWriter, r *http.Request) error {
	view.Thing().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) App(w http.ResponseWriter, r *http.Request) error {
	view.AppView().Render(r.Context(), w)
	return nil
}
