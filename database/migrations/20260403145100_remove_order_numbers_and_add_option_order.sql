-- migrate:up
-- Hapus kolom order_number yang tidak lagi diperlukan secara manual
ALTER TABLE quiz_questions DROP COLUMN IF EXISTS order_number;
ALTER TABLE quiz_options DROP COLUMN IF EXISTS order_number;

-- Tambah kolom option_order untuk menyimpan urutan pilihan jawaban yang diacak per student
-- Format JSON: {"question_id": [opt_id1, opt_id2, ...]}
ALTER TABLE quiz_attempts ADD COLUMN option_order JSONB;

-- migrate:down
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS option_order;
ALTER TABLE quiz_questions ADD COLUMN order_number INTEGER DEFAULT 0;
ALTER TABLE quiz_options ADD COLUMN order_number INTEGER DEFAULT 0;
