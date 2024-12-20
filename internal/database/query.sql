-- name: CreateUser :execresult
INSERT INTO users (email, names, password)
VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

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
SELECT ci.id, p.name, ci.quantity, p.price
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1;
