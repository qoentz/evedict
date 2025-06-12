package handler

import (
	"fmt"
	"github.com/qoentz/evedict/internal/service"
	"github.com/qoentz/evedict/internal/view/component"
	"net/http"
	"strings"
)

func Contact(mailService *service.MailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		subject := strings.TrimSpace(r.FormValue("subject"))
		message := strings.TrimSpace(r.FormValue("message"))

		if errors := mailService.ValidateContactForm(name, email, subject, message); len(errors) > 0 {
			errorMsg := strings.Join(errors, ", ")
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}

		if err := mailService.SendContactEmail(name, email, subject, message); err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err := component.SubmitButton(true).Render(r.Context(), w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
