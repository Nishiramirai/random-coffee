CREATE TABLE users (
    telegram_id   BIGINT PRIMARY KEY,
    username      TEXT,
    full_name     TEXT NOT NULL DEFAULT '',
    about         TEXT NOT NULL DEFAULT '',
    state         TEXT NOT NULL DEFAULT 'INITIAL',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rounds (
    id                 SERIAL PRIMARY KEY,
    started_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    participants_count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE matches (
    id          SERIAL PRIMARY KEY,
    round_id    INTEGER NOT NULL REFERENCES rounds(id),
    user1_id    BIGINT  NOT NULL REFERENCES users(telegram_id),
    user2_id    BIGINT  NOT NULL REFERENCES users(telegram_id),
    feedback_u1 TEXT,
    feedback_u2 TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_matches_round ON matches (round_id);
CREATE INDEX idx_matches_users ON matches (user1_id, user2_id);
CREATE INDEX idx_users_state  ON users (state);
