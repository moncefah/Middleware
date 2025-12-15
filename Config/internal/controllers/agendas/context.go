package agendas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
)

func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agendaId, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			body, status := helpers.RespondError(&models.ErrorUnprocessableEntity{
				Message: fmt.Sprintf("cannot parse id (%s) as UUID", chi.URLParam(r, "id"))})

			w.WriteHeader(status)
			if body != nil {
				_, _ = w.Write(body)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "agendaId", agendaId) // We fill context with a Key-valued variable
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
