package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplPath := filepath.Join("templates", "home.gohtml")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:mango@notmango.dev\">mango@notmango.dev</a>.</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	<h1>Welcome to the FAQ Page!</h1>
	<div>
		<ul>
			<p><b>Q:</b> Is there a free version?</p>
			<p><b>A:</b> Yes! We offer a free trial for 30 days on any paid plans.</p>
				</br>
			<p><b>Q:</b> What are your support hours?</p>
			<p><b>A:</b> We have support staff answering emails 24/7, though response times may be a bit slower on weekends</p>
				</br>
			<p><b>Q:</b> How do I contact support?</p>
			<p><b>A:</b> Email us - <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a></p>
		</ul>
	</div>
	`)
}

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "galleryID")
	w.Write([]byte(fmt.Sprintf("This is gallery id: %v", galleryID)))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.With(middleware.Logger).Get("/gallery/{galleryID}", galleriesHandler) // This adds Chi logging only to the gallery route
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
