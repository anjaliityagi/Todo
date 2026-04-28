BEGIN;


ALTER TABLE if exists users ADD COLUMN if not exists suspended_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE if exists todos RENAME COLUMN complete TO isCompleted;

CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE if not exists user_roles (
 user_id UUID NOT NULL REFERENCES users(id),
    role user_role NOT NULL,
 created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    archived_at TIMESTAMP WITH TIME ZONE


);

COMMIT; 