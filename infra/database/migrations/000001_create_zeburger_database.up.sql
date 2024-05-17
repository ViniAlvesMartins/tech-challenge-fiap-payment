CREATE SCHEMA IF NOT EXISTS ze_burguer;

CREATE TYPE ze_burguer.status_order AS ENUM ('AWAITING_PAYMENT', 'RECEIVED', 'PREPARING', 'READY', 'FINISHED');

-- Create clients table
CREATE TABLE IF NOT EXISTS ze_burguer.clients (
    "id" BIGSERIAL NOT NULL,
    "cpf" BIGINT NOT NULL UNIQUE,
    "name" VARCHAR(55) NOT NULL,
    "email" VARCHAR(55) NOT NULL UNIQUE,
    CONSTRAINT "PK_Clients" PRIMARY KEY ("id")
); 
-- Sample data clients
INSERT INTO ze_burguer.clients( "cpf", "name", "email")
VALUES (14235896700, 'dbmussarelo', 'dbmussarelo@emailo.com');

INSERT INTO ze_burguer.clients ( "cpf", "name", "email")
VALUES (14235898700, 'dbcalabresso', 'dbcalabresso@emailo.com');

INSERT INTO ze_burguer.clients ( "cpf", "name", "email")
VALUES (18535898700, 'dbtroncudo', 'dbtroncudo@emailo.com');

INSERT INTO ze_burguer.clients ( "cpf", "name", "email")
VALUES (18589898700, 'dbcasosbahio', 'dbcasosbahio@emailo.com');

INSERT INTO ze_burguer.clients ( "cpf", "name", "email")
VALUES (18589822200, 'dbludmilo', 'dbludmilo@emailo.com');

INSERT INTO ze_burguer.clients ( "cpf", "name", "email")
VALUES (18589822400, 'dbdelicio', 'dbdelicio@emailo.com');


-- Create orders table
CREATE TABLE IF NOT EXISTS ze_burguer.orders (
    "id" BIGSERIAL NOT NULL,
    "client_id" INT NULL,
    "status_order" ze_burguer.status_order NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "amount" FLOAT NOT NULL,
    CONSTRAINT "PK_order" PRIMARY KEY ("id"),
    CONSTRAINT "FK_client" FOREIGN KEY ("client_id") REFERENCES ze_burguer.clients(id)
); 
-- Sample data orders
INSERT INTO ze_burguer.orders( "client_id", "status_order", "created_at", "amount")
VALUES (1, 'AWAITING_PAYMENT', '2023-10-13 11:30:30', 17.51);
INSERT INTO ze_burguer.orders( "client_id", "status_order", "created_at", "amount")
VALUES (2, 'PREPARING', '2023-10-13 11:31:30', 20.50);
INSERT INTO ze_burguer.orders( "client_id", "status_order", "created_at", "amount")
VALUES (3, 'RECEIVED', '2023-10-13 11:32:30', 15);
INSERT INTO ze_burguer.orders( "client_id", "status_order", "created_at", "amount")
VALUES (4, 'FINISHED', '2023-10-13 11:33:30', 17);

-- Create categories table
CREATE TABLE IF NOT EXISTS ze_burguer.categories (
    "id" INT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    CONSTRAINT "PK_category" PRIMARY KEY ("id")
); 
-- Sample data category
INSERT INTO ze_burguer.categories( "id", "name")
VALUES (1, 'lanche');
INSERT INTO ze_burguer.categories( "id", "name")
VALUES (2, 'bebida');
INSERT INTO ze_burguer.categories( "id", "name")
VALUES (3, 'acompanhamento');
INSERT INTO ze_burguer.categories( "id", "name")
VALUES (4, 'sobremesa');

-- Create category product table
CREATE TABLE IF NOT EXISTS ze_burguer.products (
    "id" BIGSERIAL NOT NULL,
    "category_id" INT NOT NULL,
    "name_product" VARCHAR(255) NOT NULL,
    "price" FLOAT NOT NULL,
    "description" VARCHAR(300) NOT NULL,
    "active" BOOLEAN NOT NULL,
    CONSTRAINT "PK_products" PRIMARY KEY ("id"),
    CONSTRAINT "FK_category" FOREIGN KEY ("category_id") REFERENCES ze_burguer.categories(id)
); 
-- Sample data products
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (1, 'x-salada', 10, 'hamburgão com salada', true);
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (2, 'suco de laranja', 2.50, 'gelado', true);
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (3, 'batata frita', 5, 'quente', true);
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (1, 'x-bacon', 10, 'hamburgão com bacon', true);
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (2, 'suco de limão', 3.50, 'gelado', true);
INSERT INTO ze_burguer.products( "category_id", "name_product", "price", "description", "active")
VALUES (3, 'nuggets', 7, 'quente', true);

-- Create orders products.
CREATE TABLE IF NOT EXISTS ze_burguer.orders_products (
    "id" BIGSERIAL NOT NULL,
    "order_id" INT NOT NULL,
    "product_id" INT NOT NULL,
    CONSTRAINT "FK_id_order" FOREIGN KEY ("order_id") REFERENCES  ze_burguer.orders(id),
    CONSTRAINT "FK_id_product" FOREIGN KEY ("product_id") REFERENCES ze_burguer.products(id)
); 