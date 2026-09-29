package renders

import (
	"log"
	"net/http"
	"text/template"
)

const PathToBaseErr = "templates/errs/base_error.html"
const PathToFolderErr = "templates/errs/"

func RenderError(w http.ResponseWriter, page string, status int, data any) {
	tmpl, err := template.ParseFiles(
		PathToBaseErr,
		PathToFolderErr+page+".html",
	)
	if err != nil {
		// Если даже шаблон ошибки не парсится — отдаём голый текст
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(status) // статус 404/500/403
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.ExecuteTemplate(w, "base_error", data)

	if err != nil {
		// + logger
		log.Printf("error template: %v", err)
	}
}
