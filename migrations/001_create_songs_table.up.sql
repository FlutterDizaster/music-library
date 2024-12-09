BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS songs (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title TEXT NOT NULL,
    band TEXT NOT NULL,
    release_date DATE NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL,
    deleted BOOLEAN DEFAULT false
);

COMMIT;
