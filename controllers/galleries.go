package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"mime/multipart"

	"github.com/enlistedmango/lenslocked/models"
	"github.com/enlistedmango/lenslocked/services"
	"github.com/enlistedmango/lenslocked/views"
	"github.com/go-chi/chi"
)

type Galleries struct {
	Templates struct {
		New   views.Template
		Show  views.Template
		Index views.Template
		Edit  views.Template
	}
	GalleryService    *models.GalleryService
	FiveManageService *services.FiveManageService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: "/galleries/new",
			User:         views.ConvertToViewUser(models.UserFromContext(r.Context())),
		},
		Title: "New Gallery",
	}
	g.Templates.New.Execute(w, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	user := models.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		UserID: user.ID,
		Title:  r.FormValue("title"),
	}

	err := g.GalleryService.Create(&gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Redirect to the gallery page where they can add images
	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", gallery.ID), http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := models.UserFromContext(r.Context())
	galleries, err := g.GalleryService.GetByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: "/galleries",
			User:         views.ConvertToViewUser(user),
		},
		Galleries: galleries,
		Title:     "My Galleries",
	}
	g.Templates.Index.Execute(w, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: fmt.Sprintf("/galleries/%d", gallery.ID),
			User:         views.ConvertToViewUser(models.UserFromContext(r.Context())),
		},
		Gallery: gallery,
		Title:   gallery.Title,
	}
	g.Templates.Show.Execute(w, data)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Debug - Upload request path: %s\n", r.URL.Path)
	fmt.Printf("Debug - Request Method: %s\n", r.Method)
	fmt.Printf("Debug - Content Type: %s\n", r.Header.Get("Content-Type"))

	// Use the existing galleryByID helper instead of parsing ID manually
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// Parse the multipart form with a reasonable size limit (e.g., 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Printf("Debug - Error parsing multipart form: %v\n", err)
		http.Error(w, "Error parsing upload", http.StatusBadRequest)
		return
	}

	// Debug the entire form
	fmt.Printf("Debug - Form Values: %+v\n", r.Form)
	fmt.Printf("Debug - MultipartForm: %+v\n", r.MultipartForm)
	fmt.Printf("Debug - MultipartForm File: %+v\n", r.MultipartForm.File)

	// List all available keys in the form
	for key := range r.MultipartForm.File {
		fmt.Printf("Debug - Available file key: %s\n", key)
	}

	// Try both "image" and "images"
	var file multipart.File
	var header *multipart.FileHeader

	// Try "image" first
	file, header, err = r.FormFile("image")
	if err != nil {
		fmt.Printf("Debug - Error getting 'image': %v\n", err)
		// Try "images" as fallback
		file, header, err = r.FormFile("images")
		if err != nil {
			fmt.Printf("Debug - Error getting 'images': %v\n", err)
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			return
		}
	}
	defer file.Close()

	fmt.Printf("Debug - Successfully got file: %s, size: %d\n", header.Filename, header.Size)

	// Add metadata for FiveManage
	metadata := map[string]interface{}{
		"gallery_id": gallery.ID,
		"user_id":    gallery.UserID,
		"filename":   header.Filename,
	}

	// Upload to FiveManage
	fmResp, err := g.FiveManageService.UploadImage(file, metadata)
	if err != nil {
		fmt.Printf("Debug - FiveManage upload error: %v\n", err)
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	// Save image info to database
	err = g.GalleryService.AddImage(gallery.ID, fmResp.URL, fmResp.ID)
	if err != nil {
		fmt.Printf("Debug - Error saving to database: %v\n", err)
		http.Error(w, "Error saving image", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", gallery.ID), http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	// Try to get IDs from URL params first
	id := chi.URLParam(r, "id")
	imageIDStr := chi.URLParam(r, "imageID")

	// Fallback to path parsing if needed
	if id == "" || imageIDStr == "" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 5 {
			id = parts[2]
			imageIDStr = parts[4]
		}
	}

	fmt.Printf("Debug - Gallery ID: %q, Image ID: %q\n", id, imageIDStr)

	galleryID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return
	}
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Verify gallery ownership
	gallery, err := g.GalleryService.GetByID(galleryID)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	user := models.UserFromContext(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You don't have access to this gallery", http.StatusForbidden)
		return
	}

	// Delete from FiveManage first
	image, err := g.GalleryService.GetImage(imageID)
	if err != nil {
		fmt.Printf("Debug - Error getting image: %v\n", err)
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Debug - Attempting to delete image - DB ID: %d, FM ID: %s\n", image.ID, image.FMImageID)

	err = g.FiveManageService.DeleteImage(image.FMImageID)
	if err != nil {
		fmt.Printf("Debug - Error deleting from FiveManage: %v\n", err)
		http.Error(w, "Error deleting image from storage", http.StatusInternalServerError)
		return
	}

	// Then delete from database
	err = g.GalleryService.DeleteImage(imageID)
	if err != nil {
		http.Error(w, "Error deleting image record", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", galleryID), http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	data := views.TemplateData{
		Nav: views.NavigationData{
			CurrentRoute: fmt.Sprintf("/galleries/%d/edit", gallery.ID),
			User:         views.ConvertToViewUser(models.UserFromContext(r.Context())),
		},
		Gallery: gallery,
		Title:   "Edit Gallery",
	}
	g.Templates.Edit.Execute(w, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", gallery.ID), http.StatusFound)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// Delete all images from FiveManage first
	images, err := g.GalleryService.GetGalleryImages(gallery.ID)
	if err != nil {
		http.Error(w, "Error fetching gallery images", http.StatusInternalServerError)
		return
	}

	for _, img := range images {
		err = g.FiveManageService.DeleteImage(img.FMImageID)
		if err != nil {
			fmt.Printf("Error deleting image %s from FiveManage: %v\n", img.FMImageID, err)
			// Continue deleting other images even if one fails
		}
	}

	// Then delete the gallery (this should cascade delete the images in the database)
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// Helper function to reduce code duplication
func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	var id int
	var err error

	// Try Chi URL parameter first
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		// Fallback to path parsing
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 3 {
			idStr = parts[2]
		}
	}

	fmt.Printf("Debug - Extracted ID from URL: %q\n", idStr)

	if idStr == "" {
		http.Error(w, "Gallery ID not found", http.StatusBadRequest)
		return nil, fmt.Errorf("gallery id not found")
	}

	id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid gallery id: %w", err)
	}

	gallery, err := g.GalleryService.GetByID(id)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return nil, fmt.Errorf("gallery not found: %w", err)
	}

	user := models.UserFromContext(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You don't have access to this gallery", http.StatusForbidden)
		return nil, fmt.Errorf("unauthorized")
	}

	return gallery, nil
}
