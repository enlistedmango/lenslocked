package controllers

import "github.com/enlistedmango/lenslocked/models"

type App struct {
	GalleryService *models.GalleryService
}

func NewApp(galleryService *models.GalleryService) *App {
	return &App{
		GalleryService: galleryService,
	}
}
