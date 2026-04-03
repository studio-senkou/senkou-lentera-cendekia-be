-- migrate:up
ALTER TABLE quiz_attempts ADD COLUMN question_ids INTEGER[];

-- migrate:down
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS question_ids;
