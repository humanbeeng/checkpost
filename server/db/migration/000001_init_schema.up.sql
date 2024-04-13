CREATE TYPE "plan" AS ENUM (
  'guest',
  'free',
  'basic',
  'pro'
);

CREATE TYPE "http_method" AS ENUM (
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

CREATE TABLE "user" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" text NOT NULL,
  "username" text UNIQUE NOT NULL,
  "plan" plan NOT NULL,
  "email" varchar(60) UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE "endpoint" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "endpoint" text UNIQUE NOT NULL,
  "user_id" bigint,
  "plan" plan NOT NULL DEFAULT 'guest',
  "created_at" timestamp DEFAULT (now()),
  "expires_at" timestamp NOT NULL,
  "is_deleted" bool DEFAULT false
);

CREATE TABLE "request" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "user_id" bigint,
  "endpoint_id" bigint NOT NULL,
  "path" text NOT NULL DEFAULT '/',
  "response_id" bigint,
  "content" text,
  "method" http_method NOT NULL,
  "source_ip" text NOT NULL,
  "content_size" int NOT NULL DEFAULT 0,
  "response_code" int,
  "headers" jsonb,
  "query_params" jsonb,
  "created_at" timestamp DEFAULT (now()),
  "expires_at" timestamp NOT NULL,
  "is_deleted" bool DEFAULT false
);

CREATE TABLE "file_attachment" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "uri" text NOT NULL,
  "endpoint_id" bigint NOT NULL,
  "user_id" bigint,
  "created_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE "response" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "user_id" bigint,
  "endpoint_id" bigint NOT NULL,
  "response_code" int NOT NULL DEFAULT 200,
  "content" text,
  "created_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE INDEX "IDX_User_CreatedAt" ON "user" ("created_at");

CREATE INDEX "IDX_User_Username" ON "user" ("username");

CREATE INDEX "IDX_User_Email" ON "user" ("email");

CREATE INDEX "IDX_Endpoint_Plan_ExpiresAt" ON "endpoint" ("plan", "expires_at");

CREATE INDEX "IDX_Endpoint_ExpiresAt" ON "endpoint" ("expires_at");

CREATE INDEX "IDX_Endpoint" ON "endpoint" ("endpoint");

CREATE INDEX "IDX_Request_UserId" ON "request" ("user_id");

CREATE INDEX "IDX_Request_ExpiresAt" ON "request" ("expires_at");

CREATE INDEX "IDX_FileAttachment_UserId" ON "file_attachment" ("user_id");

CREATE INDEX "IDX_Response_UserId" ON "response" ("user_id");

COMMENT ON COLUMN "request"."source_ip" IS 'IPv4';

ALTER TABLE "endpoint" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "request" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "request" ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");

ALTER TABLE "request" ADD FOREIGN KEY ("response_id") REFERENCES "response" ("id");

ALTER TABLE "file_attachment" ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");

ALTER TABLE "file_attachment" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "response" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "response" ADD FOREIGN KEY ("endpoint_id") REFERENCES "endpoint" ("id");