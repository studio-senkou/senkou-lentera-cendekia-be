-- migrate:up

ALTER TABLE meeting_sessions ADD COLUMN mentor_user_id INTEGER;

UPDATE meeting_sessions ms
SET mentor_user_id = m.user_id
FROM mentors m
WHERE ms.mentor_id = m.id;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM meeting_sessions WHERE mentor_user_id IS NULL) THEN
        RAISE EXCEPTION 'Some meeting_sessions could not be migrated - mentor_user_id is NULL';
    END IF;
END $$;

ALTER TABLE meeting_sessions ALTER COLUMN mentor_user_id SET NOT NULL;

ALTER TABLE meeting_sessions DROP CONSTRAINT IF EXISTS fk_mentor_mt_sessions;
DROP INDEX IF EXISTS idx_meeting_sessions_mentor_id;

ALTER TABLE meeting_sessions DROP COLUMN mentor_id;

ALTER TABLE meeting_sessions RENAME COLUMN mentor_user_id TO mentor_id;

ALTER TABLE meeting_sessions
ADD CONSTRAINT fk_mentor_user_mt_sessions
FOREIGN KEY (mentor_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_meeting_sessions_mentor_id ON meeting_sessions(mentor_id);


-- migrate:down

ALTER TABLE meeting_sessions ADD COLUMN old_mentor_id INTEGER;

UPDATE meeting_sessions ms
SET old_mentor_id = (
    SELECT m.id FROM mentors m 
    WHERE m.user_id = ms.mentor_id 
    LIMIT 1
);

ALTER TABLE meeting_sessions DROP CONSTRAINT IF EXISTS fk_mentor_user_mt_sessions;
DROP INDEX IF EXISTS idx_meeting_sessions_mentor_id;

ALTER TABLE meeting_sessions DROP COLUMN mentor_id;

ALTER TABLE meeting_sessions RENAME COLUMN old_mentor_id TO mentor_id;

ALTER TABLE meeting_sessions
ADD CONSTRAINT fk_mentor_mt_sessions
FOREIGN KEY (mentor_id) REFERENCES mentors(id) ON DELETE CASCADE;

CREATE INDEX idx_meeting_sessions_mentor_id ON meeting_sessions(mentor_id);
