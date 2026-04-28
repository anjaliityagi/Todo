BEGIN;
ALTER TABLE todos RENAME COLUMN isCompleted TO is_completed;

COMMIT ;