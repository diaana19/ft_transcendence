CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id			UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name		TEXT NOT NULL,
    username	TEXT NOT NULL UNIQUE,
    password	TEXT NOT NULL,
    email		TEXT NOT NULL UNIQUE,
    bio			TEXT,
    wallpaper	TEXT,
    avatar		TEXT,
    created_at	TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at	TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at	TIMESTAMP
);