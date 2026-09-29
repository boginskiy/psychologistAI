package main

import (
	"html/template"
	"log"
	"net/http"
)

func render(w http.ResponseWriter, page string, data any) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("template error: %v", err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("execute error: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		render(w, "index", nil)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Скоро", http.StatusNotImplemented)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
