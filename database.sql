DROP DATABASE IF EXISTS cachacadb;
CREATE DATABASE IF NOT EXISTS cachacadb;
USE cachacadb;

-- Tabela Users
CREATE TABLE users
(
    id       INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email    VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone    VARCHAR(20),
    is_adm   BOOLEAN      NOT NULL DEFAULT FALSE
);

-- Tabela Products
CREATE TABLE products
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL
);

CREATE TABLE products_photos
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT          NOT NULL,
    filename   VARCHAR(255) NOT NULl,
    CONSTRAINT fk_product_photo FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE reviews
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    product_id  INT NOT NULL,
    user_id     INT NOT NULL,
    description VARCHAR(255),
    stars       INT      DEFAULT 0,
    review_date DATETIME DEFAULT NOW(),

    CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_review_product FOREIGN KEY (product_id) REFERENCES products (id)
);

-- Tabela Orders
CREATE TABLE orders
(
    id      INT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT         NOT NULL,
    status  VARCHAR(50) NOT NULL,
    date    DATETIME    NOT NULL,
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Tabela de associação para relacionamento muitos-para-muitos entre Orders e Products
CREATE TABLE order_products
(
    order_id   INT NOT NULL,
    product_id INT NOT NULL,
    PRIMARY KEY (order_id, product_id),
    CONSTRAINT fk_order_products_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_order_products_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);
