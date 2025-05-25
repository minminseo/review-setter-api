CREATE TYPE theme_color_enum AS ENUM ('dark', 'light');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    timezone VARCHAR(64) NOT NULL,
    theme_color theme_color_enum NOT NULL,
    language VARCHAR(5) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)