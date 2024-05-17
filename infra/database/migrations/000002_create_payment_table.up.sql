CREATE TYPE ze_burguer.payment_type AS ENUM ('CREDIT', 'DEBIT', 'CASH', 'PIX', 'QRCODE');

CREATE TYPE ze_burguer.payment_status AS ENUM ('PENDING', 'CONFIRMED', 'CANCELED');

CREATE TABLE IF NOT EXISTS ze_burguer.payments (
    "id" BIGSERIAL NOT NULL,
    "order_id" INT NOT NULL,
    "type" ze_burguer.payment_type NOT NULL,
    "status" ze_burguer.payment_status  NOT NULL,
    "amount" FLOAT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL,
    "deleted_at" TIMESTAMP NULL,
    CONSTRAINT "PK_Payments" PRIMARY KEY ("id"),
    CONSTRAINT "FK_Orders" FOREIGN KEY ("order_id") REFERENCES ze_burguer.orders(id)
);