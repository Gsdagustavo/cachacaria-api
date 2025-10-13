DROP
    DATABASE IF EXISTS cachacadb;

CREATE
    DATABASE IF NOT EXISTS cachacadb;

USE
    cachacadb;

-- Tabela Users
CREATE TABLE users
(
    id          INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uuid        VARCHAR(256) NOT NULL,
    email       VARCHAR(100) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    phone       VARCHAR(20),
    is_adm      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabela Products
CREATE TABLE products
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE products_photos
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT          NOT NULL,
    filename   VARCHAR(255) NOT NULl,
    CONSTRAINT fk_product_photo FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE carts
(
    id          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id     INT NOT NULL UNIQUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE carts_products
(
    id          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    cart_id     INT NOT NULL,
    product_id  INT NOT NULL,
    quantity    INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_carts_product FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE
);


CREATE TABLE reviews
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    product_id  INT NOT NULL,
    user_id     INT NOT NULL,
    description VARCHAR(255),
    stars       INT       DEFAULT 0,
    review_date DATETIME  DEFAULT NOW(),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_review_product FOREIGN KEY (product_id) REFERENCES products (id)
);
