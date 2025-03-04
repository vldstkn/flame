-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_gender AS ENUM ('male', 'female');
CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    gender user_gender NOT NULL,
    looking_for user_gender NOT NULL,
    name TEXT NOT NULL,
    bio TEXT
);
CREATE TABLE user_photos(
    id BIGSERIAL PRIMARY KEY,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    is_main boolean DEFAULT TRUE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_photos CASCADE;
DROP TABLE users CASCADE;
DROP TYPE user_gender CASCADE;
-- +goose StatementEnd
