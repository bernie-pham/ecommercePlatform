CREATE TABLE "notifications" (
    "id" bigserial PRIMARY KEY,
    "message" varchar NOT NULL,
    "recipient_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "title" varchar NOT NULL
);

ALTER TABLE "notifications" ADD FOREIGN KEY ("recipient_id") REFERENCES "users" ("id");