DROP DATABASE IF EXISTS cachacadb;
CREATE DATABASE IF NOT EXISTS cachacadb;
USE cachacadb;

-- Tabela Users
CREATE TABLE users (
                       id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       phone VARCHAR(20),
                       is_adm BOOLEAN NOT NULL DEFAULT FALSE
);

-- Tabela Products
CREATE TABLE products (
                          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          price FLOAT NOT NULL,
                          type VARCHAR(50),
                          origin VARCHAR(100),
                          manufacturer VARCHAR(100),
                          award VARCHAR(100)
);

-- Tabela Orders
CREATE TABLE orders (
                        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        user_id INT NOT NULL,
                        status VARCHAR(50) NOT NULL,
                        date DATETIME NOT NULL,
                        CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tabela de associação para relacionamento muitos-para-muitos entre Orders e Products
CREATE TABLE order_products (
                                order_id INT NOT NULL,
                                product_id INT NOT NULL,
                                PRIMARY KEY (order_id, product_id),
                                CONSTRAINT fk_order_products_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                                CONSTRAINT fk_order_products_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
