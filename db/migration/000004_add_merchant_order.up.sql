CREATE TABLE "merchant_order" (
    "id" bigserial PRIMARY KEY,
    "merchant_id" bigint NOT NULL,
    "total_price" float4 NOT NULL,
    "order_status" order_status NOT NULL DEFAULT 'open',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "order_id" bigint NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "order_items" ADD FOREIGN KEY ("merchant_order_id") REFERENCES "merchant_order" ("id");

ALTER TABLE "merchant_order" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "merchant_order" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");