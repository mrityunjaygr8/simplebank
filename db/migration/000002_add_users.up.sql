BEGIN;
CREATE TABLE IF NOT EXISTS "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "email" varchar UNIQUE NOT NULL
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
-- -- CREATE UNIQUE INDEX IF NOT EXISTS ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
COMMIT;
