-- +migrate Up
CREATE TABLE public.users (
    id bigserial NOT NULL PRIMARY KEY,
    username varchar(20) NOT NULL,
    password varchar NOT NULL,
    full_name varchar NULL,
    role varchar NOT NULL,
    is_active bool NOT NULL DEFAULT FALSE,
    refresh_token varchar NULL,
    created_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_unique UNIQUE (username)
);

INSERT INTO public.users (username, password, full_name, role, is_active) VALUES ('admin', '$2a$04$LiSuUvol8QlO76ePndhH5OzSc6vdpbewyQbFWSdyHgn3q3xTfXhnG', 'Admin', 'admin', TRUE);

-- +migrate Down
DROP TABLE public.users;

