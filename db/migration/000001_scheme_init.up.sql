CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "password_updated_at" timestamptz NOT NULL DEFAULT (now()),
  "access_level" int NOT NULL DEFAULT (1)
);


CREATE TABLE "verifications" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "is_occurpied" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "merchants" (
  "id" bigserial PRIMARY KEY,
  "country_code" int NOT NULL,
  "merchant_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" bigint NOT NULL,
  "description" varchar NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true
);

CREATE TABLE "countries" (
  "code" int PRIMARY KEY,
  "name" varchar NOT NULL,
  "continent_name" varchar NOT NULL
);

CREATE TYPE "product_status" AS ENUM (
  'out_of_stock',
  'in_stock',
  'running_low'
);
CREATE TYPE "order_status" AS ENUM (
  'open',
  'archived', 
  'canceled',
  'prepared',
  'picked',
  'on_delivery',
  'deliveried',
  'approved'
);

CREATE TABLE "products" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "merchant_id" bigint NOT NULL,
  "status" product_status,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "img_path" varchar NOT NULL
);

CREATE TABLE "product_pricing" (
  "id" bigserial PRIMARY KEY,
  "product_id" varchar NOT NULL,
  "base_price" float4 NOT NULL,
  "start_date" timestamptz NOT NULL DEFAULT (now()),
  "end_date" timestamptz NOT NULL DEFAULT (now()),
  "is_active" bool NOT NULL DEFAULT true,
  "priority" int NOT NULL DEFAULT 1
);

CREATE TABLE "product_tags" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "product_colour" (
  "id" bigserial PRIMARY KEY,
  "colour_name" varchar NOT NULL
);

CREATE TABLE "product_size" (
  "id" bigserial PRIMARY KEY,
  "size_value" varchar NOT NULL
);

CREATE TABLE "product_general_criteria" (
  "id" bigserial PRIMARY KEY,
  "criteria" varchar NOT NULL
);

CREATE TABLE "product_entry" (
  "id" bigserial PRIMARY KEY,
  "product_id" varchar NOT NULL,
  "colour_id" bigint,
  "size_id" bigint,
  "general_criteria_id" bigint,
  "quantity" int NOT NULL,
  "deal_id" bigint,
  "is_active" boolean NOT NULL DEFAULT true,
  "modified_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_tags_products" (
  "product_tags_id" varchar,
  "products_id" varchar,
  PRIMARY KEY ("product_tags_id", "products_id")
);

CREATE TABLE "order_items" (
  "order_id" bigint NOT NULL,
  "product_entry_id" bigint NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "total_price" float4 NOT NULL,
  "merchant_order_id" bigint NOT NULL
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "status" order_status NOT NULL DEFAULT 'open',
  "created_at" varchar DEFAULT (now()),
  "deal_id" bigint,
  "base_price" float4 NOT NULL,
  "discount_price" float4 NOT NULL
);

CREATE TABLE "deals" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "code" varchar NOT NULL,
  "start_date" timestamptz NOT NULL DEFAULT (now()),
  "end_date" timestamptz NOT NULL DEFAULT (now()),
  "type" varchar NOT NULL,
  "discount_rate" float4 NOT NULL,
  "merchant_id" bigint NOT NULL,
  "deal_limit" float4
);

CREATE TABLE "merchant_order" (
    "id" bigserial PRIMARY KEY,
    "merchant_id" bigint NOT NULL,
    "total_price" float4 NOT NULL,
    "order_status" order_status NOT NULL DEFAULT 'open',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "order_id" bigint NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notifications" (
    "id" bigserial PRIMARY KEY,
    "message" varchar NOT NULL,
    "recipient_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "title" varchar NOT NULL
);


CREATE TABLE "cart_item" (
  "id" bigserial PRIMARY KEY,
  "product_entry_id" bigserial NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" bigserial NOT NULL,
  "modified_at" timestamptz NOT NULL DEFAULT (now())
);



CREATE UNIQUE INDEX idx_users_email ON "users" ("email");

CREATE UNIQUE INDEX idx_deal_code ON "deals" ("code");

ALTER TABLE "merchants" ADD FOREIGN KEY ("country_code") REFERENCES "countries" ("code");

ALTER TABLE "merchants" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


ALTER TABLE "verifications" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_entry_id") REFERENCES "product_entry" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("deal_id") REFERENCES "deals" ("id"); 

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id"); 

ALTER TABLE "products" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "product_entry" ADD FOREIGN KEY ("deal_id") REFERENCES "deals" ("id");

ALTER TABLE "product_pricing" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_entry" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_entry" ADD FOREIGN KEY ("colour_id") REFERENCES "product_colour" ("id");

ALTER TABLE "product_entry" ADD FOREIGN KEY ("size_id") REFERENCES "product_size" ("id");

ALTER TABLE "product_entry" ADD FOREIGN KEY ("general_criteria_id") REFERENCES "product_general_criteria" ("id");

ALTER TABLE "product_tags_products" ADD FOREIGN KEY ("product_tags_id") REFERENCES "product_tags" ("id");

ALTER TABLE "product_tags_products" ADD FOREIGN KEY ("products_id") REFERENCES "products" ("id");

-- ALTER TABLE "deals" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("recipient_id") REFERENCES "users" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("merchant_order_id") REFERENCES "merchant_order" ("id");

ALTER TABLE "merchant_order" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "merchant_order" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

CREATE UNIQUE INDEX idx_product_tags_products ON "product_tags_products" ("product_tags_id", "products_id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("product_entry_id") REFERENCES "product_entry" ("id");

ALTER TABLE "cart_item" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");