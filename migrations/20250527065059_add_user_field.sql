-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE course
ADD COLUMN created_by INT NULL,
ADD CONSTRAINT fk_course_created_by FOREIGN KEY (created_by) REFERENCES users(id);

CREATE INDEX idx_course_created_by ON course (created_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_course_created_by;
ALTER TABLE course
DROP CONSTRAINT IF EXISTS fk_course_created_by,
DROP COLUMN IF EXISTS created_by;
-- +goose StatementEnd
