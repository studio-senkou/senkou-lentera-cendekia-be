-- migrate:up
-- Tambah kolom code sebagai kode unik kuis (8 karakter alphanumeric uppercase)
ALTER TABLE quiz_quizzes
    ADD COLUMN IF NOT EXISTS code VARCHAR(8) UNIQUE;

-- Isi kode untuk kuis yang sudah ada menggunakan random string 8 karakter
UPDATE quiz_quizzes
SET code = UPPER(SUBSTRING(MD5(RANDOM()::TEXT), 1, 8))
WHERE code IS NULL;

-- Jadikan NOT NULL setelah data lama terisi
ALTER TABLE quiz_quizzes
    ALTER COLUMN code SET NOT NULL;

-- Index untuk lookup cepat berdasarkan code
CREATE UNIQUE INDEX IF NOT EXISTS idx_quiz_quizzes_code ON quiz_quizzes(code);

-- migrate:down
DROP INDEX IF EXISTS idx_quiz_quizzes_code;
ALTER TABLE quiz_quizzes DROP COLUMN IF EXISTS code;
