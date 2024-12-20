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
INSERT INTO carts (userid)
VALUES ($1);

-- name: GetUserCart :one
SELECT * FROM carts WHERE userid = $1 LIMIT 1;

-- name: CreateCartItem :execresult
INSERT INTO cart_items (cartid, productid, quantity)
VALUES ($1, $2, $3);

-- name: GetCartItems :many
SELECT ci.id, p.name, ci.quantity, p.price
FROM cart_items ci
JOIN products p ON ci.productid = p.id
WHERE ci.cartid = $1;
