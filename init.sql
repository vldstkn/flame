CREATE EXTENSION postgis;

CREATE TYPE user_gender AS ENUM ('male', 'female');
CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    birth_date DATE,
    city TEXT,
    gender user_gender,
    bio TEXT,
    location geography(Point, 4326)
);
CREATE INDEX idx_users_location ON users USING GIST (location);
CREATE TABLE user_photos(
    id BIGSERIAL PRIMARY KEY,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    is_main boolean DEFAULT FALSE
);

CREATE TABLE preferences(
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    distance INT NOT NULL DEFAULT 100,
    age INT,
    gender user_gender,
    city TEXT
)