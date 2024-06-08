package utils

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, status int, t templ.Component) error {
	w.WriteHeader(status)

	err := t.Render(context.Background(), w)
	if err != nil {
		http.Error(w, "failed to render response template", http.StatusInternalServerError)
		return err
	}

	return nil
}
