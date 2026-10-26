CREATE SCHEMA IF NOT EXISTS taskflow;

CREATE TABLE IF NOT EXISTS taskflow.users (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(70) NOT NULL,
    phone_number VARCHAR(16) UNIQUE,

    CONSTRAINT users_version_positive
        CHECK (version > 0),

    CONSTRAINT users_full_name_not_blank
        CHECK (full_name ~ '[^[:space:]]'),

    CONSTRAINT users_phone_format
        CHECK (phone_number ~ '^\+[1-9][0-9]{6,14}$')
);

CREATE TABLE IF NOT EXISTS taskflow.tasks (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    owner_user_id INTEGER NOT NULL,

    title VARCHAR(50) NOT NULL,
    description VARCHAR(1000),

    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,

    CONSTRAINT tasks_version_positive
        CHECK (version > 0),

    CONSTRAINT tasks_title_not_blank
        CHECK (title ~ '[^[:space:]]'),


    CONSTRAINT tasks_description_not_blank
        CHECK (
            description IS NULL
            OR description ~ '[^[:space:]]'
        ),

    CONSTRAINT tasks_completion_consistent
        CHECK (
            (NOT completed AND completed_at IS NULL)
            OR
            (
                completed
                AND completed_at IS NOT NULL
                AND completed_at >= created_at
            )
        ),

    CONSTRAINT tasks_owner_user_id_fkey
        FOREIGN KEY (owner_user_id)
        REFERENCES taskflow.users(id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS tasks_owner_user_id_idx
    ON taskflow.tasks(owner_user_id);
