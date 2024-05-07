CREATE TYPE "plan" AS ENUM('free', 'basic', 'pro');

CREATE TYPE "http_method" AS ENUM(
  'get',
  'post',
  'put',
  'patch',
  'delete',
  'options',
  'head',
  'trace',
  'connect'
);

CREATE TABLE
  "user" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "name" TEXT NOT NULL,
    "avatar_url" TEXT NOT NULL DEFAULT 'https://placehold.co/32x32.png',
    "username" TEXT UNIQUE NOT NULL,
    "plan" plan NOT NULL,
    "email" VARCHAR(60) UNIQUE NOT NULL,
    "created_at" timestamptz DEFAULT (NOW()),
    "is_deleted" bool DEFAULT FALSE
  );

CREATE TABLE
  "endpoint" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "endpoint" TEXT UNIQUE NOT NULL,
    "user_id" BIGINT,
    "plan" plan NOT NULL DEFAULT 'free',
    "created_at" timestamptz DEFAULT (NOW()),
    "expires_at" timestamptz NOT NULL,
    "is_deleted" bool DEFAULT FALSE
  );

CREATE TABLE
  "request" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "uuid" TEXT NOT NULL,
    "user_id" BIGINT,
    "endpoint_id" BIGINT NOT NULL,
    "plan" plan NOT NULL DEFAULT 'free',
    "path" TEXT NOT NULL DEFAULT '/',
    "response_id" BIGINT,
    "content" TEXT,
    "method" http_method NOT NULL,
    "source_ip" TEXT NOT NULL,
    "content_size" INT NOT NULL DEFAULT 0,
    "response_code" INT,
    "headers" jsonb,
    "query_params" jsonb,
    "created_at" timestamptz DEFAULT (NOW()),
    "expires_at" timestamptz NOT NULL,
    "is_deleted" bool DEFAULT FALSE
  );

CREATE TABLE
  "file_attachment" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "uri" TEXT NOT NULL,
    "endpoint_id" BIGINT NOT NULL,
    "user_id" BIGINT,
    "created_at" timestamptz DEFAULT (NOW()),
    "is_deleted" bool DEFAULT FALSE
  );

CREATE TABLE
  "response" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "user_id" BIGINT,
    "endpoint_id" BIGINT NOT NULL,
    "response_code" INT NOT NULL DEFAULT 200,
    "content" TEXT,
    "created_at" timestamptz DEFAULT (NOW()),
    "is_deleted" bool DEFAULT FALSE
  );

CREATE INDEX "IDX_User_CreatedAt" ON "user" ("created_at");

CREATE INDEX "IDX_User_Username" ON "user" ("username");

CREATE INDEX "IDX_User_Email" ON "user" ("email");

CREATE INDEX "IDX_Endpoint_Plan_ExpiresAt" ON "endpoint" ("plan", "expires_at");

CREATE INDEX "IDX_Endpoint_ExpiresAt" ON "endpoint" ("expires_at");

CREATE INDEX "IDX_Endpoint" ON "endpoint" ("endpoint");

CREATE INDEX "IDX_Request_UserId" ON "request" ("user_id");

CREATE INDEX "IDX_Request_UUID" ON "request" ("uuid");

CREATE INDEX "IDX_Request_ExpiresAt" ON "request" ("expires_at");

CREATE INDEX "IDX_Request_Plan_ExpiresAt" ON "request" ("plan", "expires_at");

CREATE INDEX "IDX_FileAttachment_UserId" ON "file_attachment" ("user_id");

CREATE INDEX "IDX_Response_UserId" ON "response" ("user_id");

COMMENT ON COLUMN "request"."source_ip" IS 'IPv4';

ALTER TABLE "endpoint"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "request"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "request"
ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");

ALTER TABLE "request"
ADD FOREIGN KEY ("response_id") REFERENCES "response" ("id");

ALTER TABLE "file_attachment"
ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");

ALTER TABLE "file_attachment"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "response"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "response"
ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");