CREATE TABLE "cart_item" (
  "id" bigserial PRIMARY KEY,
  "product_entry_id" bigserial NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" bigserial NOT NULL,
  "modified_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "cart_item" ADD FOREIGN KEY ("product_entry_id") REFERENCES "product_entry" ("id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");