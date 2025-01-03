package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/enlistedmango/lenslocked/controllers"
	"github.com/enlistedmango/lenslocked/middleware"
	"github.com/enlistedmango/lenslocked/models"
	"github.com/enlistedmango/lenslocked/services"
	"github.com/enlistedmango/lenslocked/views"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func parseTemplate(filepath string) (views.Template, error) {
	return views.Parse(
		filepath,
		"templates/shared/header.gohtml",
		"templates/shared/footer.gohtml",
		"templates/shared/nav.gohtml",
		"templates/shared/theme-switcher.gohtml",
	)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Load all configuration
	cfg := models.LoadConfig()

	// Setup database
	db, err := models.OpenWithRetry(cfg.Postgres, 5, 5*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	fmt.Println("Connected to database!")

	// Setup our model services
	userService := &models.UserService{
		DB: db,
	}

	// Setup session store with config
	sessionSecret := []byte(cfg.Session.Key)
	sessionStore := sessions.NewCookieStore(sessionSecret)

	// Setup auth middleware
	authMiddleware := middleware.AuthMiddleware{
		Store:       sessionStore,
		UserService: userService,
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.Debug)
	r.Use(authMiddleware.SetUser)

	// Setup users controller
	usersC := controllers.Users{
		UserService:  userService,
		SessionStore: sessionStore,
	}
	usersC.Templates.New = views.Must(parseTemplate("templates/signup.gohtml"))
	usersC.Templates.SignIn = views.Must(parseTemplate("templates/signin.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Post("/signout", usersC.SignOut)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	// Create gallery service
	galleryService := &models.GalleryService{
		DB: db,
	}

	// Create FiveManage service with config
	fiveManageService := &services.FiveManageService{
		APIKey: cfg.FiveManage.APIKey,
		Debug:  cfg.FiveManage.Debug,
	}

	// Initialize App with services
	app := controllers.NewApp(galleryService)

	// Setup static routes
	r.Get("/", controllers.StaticHandler(app, views.Must(parseTemplate("templates/home.gohtml")), "/"))
	r.Get("/contact", controllers.StaticHandler(app, views.Must(parseTemplate("templates/contact.gohtml")), "/contact"))
	r.Get("/faq", controllers.StaticHandler(app, views.Must(parseTemplate("templates/faq.gohtml")), "/faq"))

	// Setup galleries controller
	galleriesC := controllers.Galleries{
		GalleryService:    galleryService,
		FiveManageService: fiveManageService,
	}

	// Parse templates with their required layouts
	galleriesC.Templates.New = views.Must(parseTemplate("templates/galleries/new.gohtml"))
	galleriesC.Templates.Index = views.Must(parseTemplate("templates/galleries/index.gohtml"))
	galleriesC.Templates.Show = views.Must(parseTemplate("templates/galleries/show.gohtml"))
	galleriesC.Templates.Edit = views.Must(parseTemplate("templates/galleries/edit.gohtml"))

	// Create a router group for authenticated routes
	userMw := middleware.RequireUser{}

	// Create a subrouter for galleries
	galleryRouter := chi.NewRouter()
	galleryRouter.Use(userMw.Apply)
	galleryRouter.Use(middleware.Debug)

	// List galleries
	galleryRouter.Get("/", galleriesC.Index)
	galleryRouter.Get("/new", galleriesC.New)
	galleryRouter.Post("/", galleriesC.Create)

	// Create a subrouter for gallery-specific routes
	galleryRouter.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Use(middleware.Debug)
		r.Get("/", galleriesC.Show)
		r.Get("/edit", galleriesC.Edit)
		r.Post("/", galleriesC.Update)
		r.Post("/delete", galleriesC.Delete)

		// Image routes
		r.Group(func(r chi.Router) {
			r.Post("/images", galleriesC.UploadImage)
			r.Post("/images/{imageID:[0-9]+}", galleriesC.DeleteImage)
		})
	})

	// Mount the gallery router
	r.Mount("/galleries", galleryRouter)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Add this after setting up all routes but before starting the server
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	// Add health check endpoint BEFORE starting the server
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			fmt.Printf("Health check failed: %v\n", err)
			http.Error(w, fmt.Sprintf("Database connection error: %v", err), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","database":"connected"}`)
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback for local development
	}

	fmt.Printf("Starting the server on :%s...\n", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		return
	}
	fmt.Println("Server failed to start!")
}
