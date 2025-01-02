package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Gallery struct {
	ID        int
	UserID    int
	Title     string
	Images    []GalleryImage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GalleryImage struct {
	ID        int
	GalleryID int
	URL       string
	FMImageID string // FiveManage Image ID for deletion
	CreatedAt time.Time
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(gallery *Gallery) error {
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (user_id, title)
		VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID, &gallery.CreatedAt, &gallery.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (gs *GalleryService) AddImage(galleryID int, url, fmImageID string) error {
	fmt.Printf("Debug - AddImage - Gallery ID: %d, URL: %s, FM ID: %s\n", galleryID, url, fmImageID)

	var imageID int
	err := gs.DB.QueryRow(`
		INSERT INTO gallery_images (gallery_id, url, fm_image_id)
		VALUES ($1, $2, $3)
		RETURNING id`,
		galleryID, url, fmImageID).Scan(&imageID)

	if err != nil {
		return fmt.Errorf("error inserting image: %w", err)
	}

	fmt.Printf("Debug - Image inserted with ID: %d\n", imageID)
	return nil
}

func (gs *GalleryService) GetByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title, created_at, updated_at
		FROM galleries
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying galleries: %w", err)
	}
	defer rows.Close()

	var galleries []Gallery
	for rows.Next() {
		var gallery Gallery
		gallery.UserID = userID
		err := rows.Scan(
			&gallery.ID,
			&gallery.Title,
			&gallery.CreatedAt,
			&gallery.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning gallery: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating galleries: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) GetByID(id int) (*Gallery, error) {
	fmt.Printf("Fetching gallery with ID: %d\n", id)

	gallery := Gallery{
		ID: id,
	}
	row := gs.DB.QueryRow(`
		SELECT user_id, title, created_at, updated_at
		FROM galleries WHERE id = $1`, id)
	err := row.Scan(
		&gallery.UserID,
		&gallery.Title,
		&gallery.CreatedAt,
		&gallery.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("gallery not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching gallery: %w", err)
	}

	fmt.Printf("Found gallery: %+v\n", gallery)

	// Fetch images for this gallery
	rows, err := gs.DB.Query(`
		SELECT id, url, fm_image_id, created_at
		FROM gallery_images
		WHERE gallery_id = $1
		ORDER BY created_at DESC`, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching gallery images: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var image GalleryImage
		image.GalleryID = id
		err := rows.Scan(
			&image.ID,
			&image.URL,
			&image.FMImageID,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning gallery image: %w", err)
		}
		gallery.Images = append(gallery.Images, image)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating gallery images: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) GetImage(id int) (*GalleryImage, error) {
	var image GalleryImage
	err := gs.DB.QueryRow(`
		SELECT id, gallery_id, url, fm_image_id, created_at
		FROM gallery_images WHERE id = $1`, id).Scan(
		&image.ID, &image.GalleryID, &image.URL, &image.FMImageID, &image.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %w", err)
	}
	fmt.Printf("Debug - GetImage - ID: %d, Gallery ID: %d, FM ID: %s\n",
		image.ID, image.GalleryID, image.FMImageID)
	return &image, nil
}

func (gs *GalleryService) DeleteImage(id int) error {
	_, err := gs.DB.Exec(`DELETE FROM gallery_images WHERE id = $1`, id)
	return err
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	row := gs.DB.QueryRow(`
		UPDATE galleries
		SET title = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at`,
		gallery.Title, gallery.ID)
	err := row.Scan(&gallery.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error updating gallery: %w", err)
	}
	return nil
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, id)
	return err
}

func (service *GalleryService) GetGalleryImages(galleryID int) ([]GalleryImage, error) {
	rows, err := service.DB.Query(`
		SELECT id, gallery_id, url, fm_image_id, created_at
		FROM gallery_images
		WHERE gallery_id = $1;
	`, galleryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []GalleryImage
	for rows.Next() {
		var image GalleryImage
		err := rows.Scan(
			&image.ID,
			&image.GalleryID,
			&image.URL,
			&image.FMImageID,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}
