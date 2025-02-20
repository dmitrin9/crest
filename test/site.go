package main

import (
	"html/template"
	"log"
	"net/http"
)

func run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("views/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Initialized.")

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Executed.")

	})

	http.HandleFunc("/ActuallyExists", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("views/ActuallyExists.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Initialized.")

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Executed.")

	})

	log.Println("serving 8080...")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	log.Fatal(run())
}
