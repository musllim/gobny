-- name: CreateUser :execresult
INSERT INTO users (email, names, password)
VALUES ($1, $2, $3);
