-- +goose Up
CREATE TABLE galleries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE gallery_images (
    id SERIAL PRIMARY KEY,
    gallery_id INTEGER NOT NULL REFERENCES galleries(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    fm_image_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE gallery_images;
DROP TABLE galleries; 
