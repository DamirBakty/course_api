-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE chapter
    ADD COLUMN created_by BIGINT NULL,
ADD CONSTRAINT fk_chapter_created_by FOREIGN KEY (created_by) REFERENCES users(id);

CREATE INDEX idx_chapter_created_by ON chapter (created_by);

ALTER TABLE lesson
    ADD COLUMN created_by BIGINT NULL,
ADD CONSTRAINT fk_lesson_created_by FOREIGN KEY (created_by) REFERENCES users(id);

CREATE INDEX idx_lesson_created_by ON lesson (created_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_lesson_created_by;
ALTER TABLE lesson
DROP CONSTRAINT IF EXISTS fk_lesson_created_by,
DROP COLUMN IF EXISTS created_by;

DROP INDEX IF EXISTS idx_chapter_created_by;
ALTER TABLE chapter
DROP CONSTRAINT IF EXISTS fk_chapter_created_by,
DROP COLUMN IF EXISTS created_by;
-- +goose StatementEnd
