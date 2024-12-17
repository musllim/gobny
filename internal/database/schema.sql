CREATE TABLE users (
    id SERIAL PRIMARY KEY,  
    email VARCHAR(255) UNIQUE,
    names VARCHAR(255),
    password VARCHAR(255),
    isVerified BOOLEAN DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,  
    name VARCHAR(255),
    price DECIMAL(10, 2),
    count INT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE product_images (
--     id SERIAL PRIMARY KEY,  
--     product_id INT REFERENCES products(id) ON DELETE CASCADE,  
--     image_url VARCHAR(255)
-- );

-- CREATE TABLE carts (
--     id SERIAL PRIMARY KEY,  
--     user_id INT REFERENCES users(id) ON DELETE CASCADE,  
--     createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE cart_items (
--     id SERIAL PRIMARY KEY,  
--     cart_id INT REFERENCES carts(id) ON DELETE CASCADE,  
--     product_id INT REFERENCES products(id) ON DELETE CASCADE,  
--     quantity INT,
--     createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
