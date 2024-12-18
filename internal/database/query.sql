-- name: CreateUser :execresult
INSERT INTO users (email, names, password)
VALUES ($1, $2, $3);

-- name: CreateProduct :execresult
INSERT INTO products (name, price, count)
VALUES ($1, $2, $3);

-- name: GetProducts :many
SELECT * FROM products;

-- name: CreateUserCart :execresult
INSERT INTO carts (user_id)
VALUES ($1);

-- name: GetUserCart :many
SELECT * FROM carts WHERE user_id = $1;

-- name: CreateCartItem :execresult
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3);

-- name: GetCartItems :many
SELECT * FROM cart_items WHERE cart_id = $1;
