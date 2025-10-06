-- +goose Up

CREATE TABLE USERS (
                       ID BIGSERIAL PRIMARY KEY,
                       NAME TEXT NOT NULL,
                       EMAIL TEXT NOT NULL
);

CREATE TABLE EVENTS (
                        ID BIGSERIAL PRIMARY KEY,
                        TITLE TEXT NOT NULL,
                        DESCRIPTION TEXT,
                        START_TIME TIMESTAMP NOT NULL,
                        END_TIME TIMESTAMP NOT NULL,
                        USER_ID INT NOT NULL,
                        NOTIFY_PERIOD BIGINT DEFAULT 0, --nanoseconds
                        CONSTRAINT fk_events_users FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_events__user_id ON events(user_id);

-- +goose Down

DROP TABLE EVENTS;
DROP TABLE USERS;