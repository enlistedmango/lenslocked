package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:mango@notmango.dev\">mango@notmango.dev</a>.</p>")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
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
	r.Get("/gallery/{galleryID}", galleriesHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
