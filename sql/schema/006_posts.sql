-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(250) NOT NULL,    
    url VARCHAR(500) NOT NULL UNIQUE,
    description VARCHAR(1000),
    published_at TIMESTAMP,
    feed_id UUID REFERENCES feeds (id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE posts;