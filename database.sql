DROP DATABASE IF EXISTS cachacadb;

CREATE DATABASE IF NOT EXISTS cachacadb;

USE cachacadb;

CREATE TABLE users
(
    id          INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uuid        VARCHAR(256) NOT NULL,
    email       VARCHAR(100) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    phone       VARCHAR(20),
    is_adm      BOOLEAN      NOT NULL DEFAULT FALSE,
    status_code TINYINT(1)   NOT NULL DEFAULT 0,
    created_at  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE products
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL,
    status_code TINYINT(1)     NOT NULL DEFAULT 0,
    created_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE products_photos
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT          NOT NULL,
    filename   VARCHAR(255) NOT NULl,
    CONSTRAINT fk_product_photo FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE carts_products
(
    id          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id     INT NOT NULL,
    product_id  INT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_carts_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE orders
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    user_id     INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_user FOREIGN KEY (user_id) REFERENCES users (id)
);


CREATE TABLE order_items
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    order_id    INT NOT NULL,
    product_id  INT NOT NULL,
    quantity    INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_order_item_product FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO users (uuid, email, password, phone, is_adm, status_code)
VALUES (UUID(), 'admin@wilbert.com', '$2a$10$O55dgMZop3M67kLi.GV/RuQQlgNc1G.4yAnqzzzDAJZ02hBR2MVge', '47999999999',
        TRUE, 1);

INSERT INTO products (name, description, price, stock, status_code)
VALUES ('Cachaça Ouro da Serra 700ml',
        'Cachaça artesanal armazenada em barris de carvalho por 2 anos, sabor amadeirado e suave.', 59.90, 35, 1),
       ('Cachaça Prata do Vale 670ml', 'Cachaça branca tradicional, destilada em alambiques de cobre com aroma leve.',
        42.00, 50, 1),
       ('Cachaça Reserva Mineira 750ml', 'Envelhecida por 5 anos em tonéis de amburana, notas doces e encorpadas.',
        89.90, 20, 1),
       ('Cachaça Engenho Velho 600ml', 'Produção limitada, armazenada em jequitibá, sabor equilibrado.', 55.00, 40, 1),
       ('Cachaça Vale da Canela 700ml', 'Aromática e levemente adocicada, ideal para caipirinhas especiais.', 38.50, 60,
        1),
       ('Cachaça Carvalheira Premium 750ml', 'Envelhecida em carvalho europeu, corpo intenso e final marcante.', 110.00,
        15, 1),
       ('Cachaça Alambique Real 500ml', 'Cachaça de alambique com graduação suave e aroma herbal.', 29.90, 80, 1),
       ('Cachaça Ouro Nacional 700ml', 'Amadurecida por 3 anos em carvalho, equilibrada e macia.', 65.00, 30, 1),
       ('Cachaça Colonial Prata 1L', 'Branca e leve, perfeita para drinks e caipirinhas.', 35.90, 55, 1),
       ('Cachaça Jequitibá Dourada 750ml', 'Armazenada em jequitibá rosa, sabor levemente adocicado.', 58.00, 25, 1),
       ('Cachaça Amburana Premium 700ml', 'Envelhecida exclusivamente em amburana, aroma intenso e toque de baunilha.',
        72.00, 22, 1),
       ('Cachaça Flor do Engenho 600ml', 'Cachaça artesanal branca, sabor fresco e suave.', 31.90, 70, 1),
       ('Cachaça Canavial Ouro 750ml', 'Ouro, maturada em carvalho por 18 meses, sabor amadeirado.', 49.90, 45, 1),
       ('Cachaça Forte Cana 1L', 'Cachaça encorpada e aromática, ideal para coquetéis.', 28.00, 90, 1),
       ('Cachaça Diamantina Prata 700ml', 'Branca premium, dupla destilação, sabor extremamente puro.', 55.90, 35, 1),
       ('Cachaça Sabor da Roça 750ml', 'Produção rústica com notas herbais e final seco.', 39.90, 65, 1),
       ('Cachaça Castanheira 700ml', 'Armazenada em tonéis de castanheira, sabor distinto e adocicado.', 63.00, 28, 1),
       ('Cachaça Baronesa Ouro 750ml', 'Ouro, longa maturação e aroma frutado.', 79.00, 18, 1),
       ('Cachaça Cristalina 900ml', 'Branca, extremamente leve e cristalina.', 27.50, 100, 1),
       ('Cachaça Monte Belo 700ml', 'Tradicional, notas terrosas e corpo moderado.', 33.00, 75, 1),
       ('Cachaça Engenho Antigo 750ml', 'Cachaça envelhecida em bálsamo por 3 anos.', 71.00, 20, 1),
       ('Cachaça Vale Encantado 700ml', 'Cachaça premium de alambique, aroma intenso e sabor equilibrado.', 47.00, 52,
        1),
       ('Cachaça Ouro Imperial 750ml', 'Carvalho francês, 4 anos de envelhecimento.', 120.00, 12, 1),
       ('Cachaça Tapira Prata 700ml', 'Branca suave, excelente custo-benefício.', 29.50, 85, 1),
       ('Cachaça Horizonte Ouro 750ml', 'Notas de mel e baunilha, envelhecida em amburana.', 66.00, 26, 1),
       ('Cachaça Engenho da Mata 700ml', 'Produção sustentável com cana orgânica.', 54.90, 33, 1),
       ('Cachaça Dom Alambique 750ml', 'Premium artesanal, blend de carvalho e amburana.', 95.00, 17, 1),
       ('Cachaça Tradição Prata 600ml', 'Branca, aroma fresco e toque cítrico.', 26.90, 95, 1),
       ('Cachaça Bambu Ouro 750ml', 'Armazenada em tonéis de bambu, sabor suave incomparável.', 57.50, 34, 1),
       ('Cachaça Serenata 700ml', 'Cachaça premium com acabamento frutado.', 69.90, 21, 1);
