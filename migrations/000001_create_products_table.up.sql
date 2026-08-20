CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    company VARCHAR(100) NOT NULL,
    price INT NOT NULL CHECK(price >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_company ON products(company);