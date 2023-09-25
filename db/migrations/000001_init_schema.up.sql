CREATE TABLE "tenant" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar UNIQUE NOT NULL,
  "logo" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "username" varchar(50) NOT NULL,
  "email" varchar(40) NOT NULL,
  "tenant_id" uuid NOT NULL,
  "password" varchar(60) NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "logs" (
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "user_id" uuid NOT NULL,
    "description" varchar(40) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("id", "tenant_id");
ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id");
ALTER TABLE "logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");