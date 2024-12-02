CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ,
                       username VARCHAR(100) NOT NULL,
                       password VARCHAR(100) NOT NULL,
                       role VARCHAR(50) DEFAULT 'user' NOT NULL,
                       CONSTRAINT uni_users_username UNIQUE (username)
);

CREATE TABLE groups (
                        id SERIAL PRIMARY KEY,
                        created_at TIMESTAMPTZ DEFAULT NOW(),
                        updated_at TIMESTAMPTZ DEFAULT NOW(),
                        deleted_at TIMESTAMPTZ,
                        name VARCHAR(100) NOT NULL,
                        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE urls (
                      id SERIAL PRIMARY KEY,
                      created_at TIMESTAMPTZ DEFAULT NOW(),
                      updated_at TIMESTAMPTZ DEFAULT NOW(),
                      deleted_at TIMESTAMPTZ,
                      address TEXT NOT NULL,
                      group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE
);
