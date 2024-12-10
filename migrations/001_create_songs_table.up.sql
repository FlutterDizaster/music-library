BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS songs (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    band VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    lyrics TEXT NOT NULL,
    link VARCHAR(255) NOT NULL,
    deleted BOOLEAN DEFAULT false
);

COMMIT;
