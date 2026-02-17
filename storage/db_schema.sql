PRAGMA foreign_keys = ON;

CREATE TABLE
  IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY,
    "name" TEXT,
    "surname" TEXT,
    "nickname" TEXT UNIQUE,
    "created_at" TEXT,
    "updated_at" TEXT
  );

CREATE TABLE
  IF NOT EXISTS "auth" (
    "user_id" INTEGER PRIMARY KEY,
    "email" TEXT UNIQUE,
    "password_hash" TEXT,
    "created_at" TEXT,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
  );

CREATE TABLE
  IF NOT EXISTS "requests" (
    "id" INTEGER PRIMARY KEY,
    "request_text" TEXT,
    "user_id" INTEGER,
    "created_at" TEXT,
    "updated_at" TEXT,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
  );

CREATE INDEX IF NOT EXISTS "idx_requests_user_id" ON "requests" ("user_id");

CREATE TABLE
  IF NOT EXISTS "questions" (
    "id" INTEGER PRIMARY KEY,
    "request_id" INTEGER,
    "question" TEXT,
    "answer" TEXT,
    "user_id" INTEGER,
    "level" TEXT,
    "created_at" TEXT,
    "updated_at" TEXT,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("request_id") REFERENCES "requests" ("id") ON DELETE CASCADE
  );

CREATE INDEX IF NOT EXISTS "idx_questions_request_id" ON "questions" ("request_id");

CREATE INDEX IF NOT EXISTS "idx_questions_user_id" ON "questions" ("user_id");