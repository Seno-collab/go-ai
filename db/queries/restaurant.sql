-- name: CreateRestaurant :one
INSERT INTO "restaurant" (name, description, address, category, city, district, logo_url, banner_url, phone_number, website_url, email)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: GetByName :one
SELECT * FROM "restaurant" WHERE name LIKE $1;

-- name: GetById :one
SELECT * FROM "restaurant" WHERE id = $1;
