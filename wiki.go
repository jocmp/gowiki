package main

import (
	"html/template"
	"log"
	"net/http"
)	

type route string

const (
	view route = "/view/"
	save route = "/save/"
	edit route = "/edit/"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(view):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, string(edit)+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(edit):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(save):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, string(view)+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc(string(edit), editHandler)
	http.HandleFunc(string(view), viewHandler)
	http.HandleFunc(string(save), saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
