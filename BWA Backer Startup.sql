CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "occupation" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password_hash" varchar NOT NULL,
  "token" varchar,
  "avatar_file_name" varchar,
  "role" varchar,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "campaigns" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "user_id" int NOT NULL,
  "short_description" varchar NOT NULL,
  "description" text NOT NULL,
  "goal_amount" int NOT NULL,
  "current_amount" int NOT NULL DEFAULT 0,
  "backer_count" int NOT NULL DEFAULT 0,
  "perks" text NOT NULL,
  "slug" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "campaign_images" (
  "id" bigserial PRIMARY KEY,
  "campain_id" int NOT NULL,
  "file_name" varchar NOT NULL,
  "is_primary" bool NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "campain_id" int NOT NULL,
  "user_id" int NOT NULL,
  "amount" int NOT NULL,
  "code" varchar NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

ALTER TABLE "campaign" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "campaign_images" ADD FOREIGN KEY ("campain_id") REFERENCES "campaigns" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("campain_id") REFERENCES "campaigns" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
