CREATE TYPE meeting_format AS ENUM ('online', 'offline', 'any');

ALTER TABLE users
    ADD COLUMN preferred_format meeting_format NOT NULL DEFAULT 'any';
