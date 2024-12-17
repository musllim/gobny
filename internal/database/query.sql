-- name: CreateUser :execresult
INSERT INTO users (email, names, password)
VALUES ($1, $2, $3);

-- name: CreateProduct :execresult
INSERT INTO products (name, price, count)
VALUES ($1, $2, $3);

-- name: GetProducts :many
SELECT * FROM products;