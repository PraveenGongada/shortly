CREATE TABLE IF NOT EXISTS "user" (
    "id" character(36) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);

CREATE TABLE IF NOT EXISTS url (
    "id" character(36) NOT NULL PRIMARY KEY,
    "user_id" character(36) NOT NULL REFERENCES "user"(id),
    "short_url" varchar(7) NOT NULL UNIQUE,
    "long_url" TEXT NOT NULL,
    "redirects" INT NOT NULL DEFAULT 0,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);
