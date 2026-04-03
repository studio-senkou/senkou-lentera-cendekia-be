package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/studio-senkou/lentera-cendekia-be/database/facades"
)

type QuizQuiz struct {
	ID               uint       `json:"id"`
	Code             string     `json:"code"`
	Title            string     `json:"title"`
	Description      *string    `json:"description,omitempty"`
	PassingScore     int        `json:"passing_score"`
	TimeLimitMinutes *int       `json:"time_limit_minutes,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

func generateQuizCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}

type QuizQuestion struct {
	ID           uint         `json:"id"`
	QuizID       uint         `json:"quiz_id"`
	QuestionText string       `json:"question_text"`
	Options      []QuizOption `json:"options,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at"`
}

type QuizOption struct {
	ID         uint       `json:"id"`
	QuestionID uint       `json:"question_id"`
	OptionText string     `json:"option_text"`
	IsCorrect  bool       `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type QuizOptionAdmin struct {
	ID         uint       `json:"id"`
	QuestionID uint       `json:"question_id"`
	OptionText string     `json:"option_text"`
	IsCorrect  bool       `json:"is_correct"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type QuizAttempt struct {
	ID                   uint            `json:"id"`
	QuizID               uint            `json:"quiz_id"`
	UserID               uint            `json:"user_id"`
	Status               string          `json:"status"`
	Score                *float64        `json:"score"`
	QuestionIDs          pq.Int64Array   `json:"-"`
	OptionOrder          json.RawMessage `json:"-"`
	CurrentQuestionIndex int             `json:"-"`
	StartedAt            time.Time       `json:"started_at"`
	SubmittedAt          *time.Time      `json:"submitted_at"`
	ResetAt              *time.Time      `json:"reset_at,omitempty"`
	ResetBy              *uint           `json:"reset_by,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            *time.Time      `json:"updated_at"`
}

type QuizAnswer struct {
	ID         uint      `json:"id"`
	AttemptID  uint      `json:"attempt_id"`
	QuestionID uint      `json:"question_id"`
	OptionID   uint      `json:"option_id"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
}

type QuizRepository struct {
	db  facades.DBExecutor
	raw *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db, raw: db}
}

func (r *QuizRepository) WithExecutor(executor facades.DBExecutor) *QuizRepository {
	return &QuizRepository{db: executor, raw: r.raw}
}

func (r *QuizRepository) GetActiveQuizWithQuestions(quizID uint, questionIDs []int64) (*QuizQuiz, []QuizQuestion, error) {
	quizQuery := `
		SELECT id, code, title, description, passing_score, time_limit_minutes, is_active, created_at, updated_at
		FROM quiz_quizzes
		WHERE id = $1 AND is_active = TRUE AND deleted_at IS NULL
	`
	quiz := new(QuizQuiz)
	err := r.db.QueryRow(quizQuery, quizID).Scan(
		&quiz.ID, &quiz.Code, &quiz.Title, &quiz.Description, &quiz.PassingScore,
		&quiz.TimeLimitMinutes, &quiz.IsActive, &quiz.CreatedAt, &quiz.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var qRows *sql.Rows
	if len(questionIDs) > 0 {
		questionsQuery := `
			SELECT id, quiz_id, question_text, created_at, updated_at
			FROM quiz_questions
			WHERE id = ANY($1)
			ORDER BY array_position($1, id)
		`
		qRows, err = r.db.Query(questionsQuery, pq.Array(questionIDs))
	} else {
		questionsQuery := `
			SELECT id, quiz_id, question_text, created_at, updated_at
			FROM quiz_questions
			WHERE quiz_id = $1
			ORDER BY id ASC
		`
		qRows, err = r.db.Query(questionsQuery, quizID)
	}

	if err != nil {
		return nil, nil, err
	}
	defer qRows.Close()

	questions := make([]QuizQuestion, 0)
	loadedIDs := make([]uint, 0)
	questionMap := make(map[uint]*QuizQuestion)

	for qRows.Next() {
		q := QuizQuestion{}
		if err := qRows.Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, nil, err
		}
		q.Options = make([]QuizOption, 0)
		questions = append(questions, q)
		loadedIDs = append(loadedIDs, q.ID)
	}

	if len(loadedIDs) > 0 {
		optionsQuery := `
			SELECT id, question_id, option_text, created_at, updated_at
			FROM quiz_options
			WHERE question_id = ANY($1)
		`
		int64IDs := make([]int64, len(loadedIDs))
		for i, id := range loadedIDs {
			int64IDs[i] = int64(id)
		}
		oRows, err := r.db.Query(optionsQuery, pq.Array(int64IDs))
		if err != nil {
			return nil, nil, err
		}
		defer oRows.Close()

		for _, q := range questions {
			qCopy := q
			questionMap[q.ID] = &qCopy
		}

		for oRows.Next() {
			opt := QuizOption{}
			if err := oRows.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText, &opt.CreatedAt, &opt.UpdatedAt); err != nil {
				return nil, nil, err
			}
			if q, ok := questionMap[opt.QuestionID]; ok {
				q.Options = append(q.Options, opt)
			}
		}

		for i, q := range questions {
			if updated, ok := questionMap[q.ID]; ok {
				questions[i] = *updated
			}
		}
	}

	return quiz, questions, nil
}

func (r *QuizRepository) GetActiveQuizWithQuestionsV2(quizID uint, attempt *QuizAttempt) (*QuizQuiz, []QuizQuestion, error) {
	quiz, questions, err := r.GetActiveQuizWithQuestions(quizID, attempt.QuestionIDs)
	if err != nil {
		return nil, nil, err
	}

	if attempt.OptionOrder != nil {
		var optionOrder map[string][]int64
		if err := json.Unmarshal(attempt.OptionOrder, &optionOrder); err == nil {
			for i, q := range questions {
				qIDStr := strconv.FormatUint(uint64(q.ID), 10)
				if order, ok := optionOrder[qIDStr]; ok {
					orderedOptions := make([]QuizOption, 0, len(q.Options))
					optMap := make(map[int64]QuizOption)
					for _, opt := range q.Options {
						optMap[int64(opt.ID)] = opt
					}
					for _, optID := range order {
						if opt, found := optMap[optID]; found {
							orderedOptions = append(orderedOptions, opt)
						}
					}
					questions[i].Options = orderedOptions
				}
			}
		}
	}

	return quiz, questions, nil
}

func (r *QuizRepository) GetQuizIDByCode(code string) (uint, error) {
	var id uint
	err := r.db.QueryRow(
		`SELECT id FROM quiz_quizzes WHERE code = $1 AND is_active = TRUE AND deleted_at IS NULL`,
		code,
	).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (r *QuizRepository) GetActiveAttempt(userID, quizID uint) (*QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, status, score, question_ids, option_order,
		       current_question_index, started_at, submitted_at, reset_at, reset_by,
		       created_at, updated_at
		FROM quiz_attempts
		WHERE user_id = $1 AND quiz_id = $2 AND status != 'reset'
		ORDER BY created_at DESC
		LIMIT 1
	`
	attempt := new(QuizAttempt)
	err := r.db.QueryRow(query, userID, quizID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.Status, &attempt.Score,
		&attempt.QuestionIDs, &attempt.OptionOrder, &attempt.CurrentQuestionIndex,
		&attempt.StartedAt, &attempt.SubmittedAt,
		&attempt.ResetAt, &attempt.ResetBy, &attempt.CreatedAt, &attempt.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return attempt, nil
}

func (r *QuizRepository) UpdateAttemptIndex(attemptID uint, index int) error {
	_, err := r.db.Exec(
		`UPDATE quiz_attempts SET current_question_index = $1 WHERE id = $2`,
		index, attemptID,
	)
	return err
}

func (r *QuizRepository) GetQuestionByAttemptIndex(attempt *QuizAttempt, index int) (*QuizQuestion, error) {
	if index < 0 || index >= len(attempt.QuestionIDs) {
		return nil, nil
	}
	qID := attempt.QuestionIDs[index]

	q := &QuizQuestion{}
	err := r.db.QueryRow(
		`SELECT id, quiz_id, question_text, created_at, updated_at FROM quiz_questions WHERE id = $1`,
		qID,
	).Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var optionIDs []int64
	if attempt.OptionOrder != nil {
		var optionOrder map[string][]int64
		if err := json.Unmarshal(attempt.OptionOrder, &optionOrder); err == nil {
			qIDStr := strconv.FormatInt(qID, 10)
			if order, ok := optionOrder[qIDStr]; ok {
				optionIDs = order
			}
		}
	}

	if len(optionIDs) > 0 {
		rows, err := r.db.Query(
			`SELECT id, question_id, option_text, created_at, updated_at FROM quiz_options WHERE id = ANY($1)`,
			pq.Array(optionIDs),
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		optMap := make(map[int64]QuizOption)
		for rows.Next() {
			var opt QuizOption
			if err := rows.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText, &opt.CreatedAt, &opt.UpdatedAt); err != nil {
				return nil, err
			}
			optMap[int64(opt.ID)] = opt
		}
		for _, oid := range optionIDs {
			if opt, ok := optMap[oid]; ok {
				q.Options = append(q.Options, opt)
			}
		}
	} else {
		rows, err := r.db.Query(
			`SELECT id, question_id, option_text, created_at, updated_at FROM quiz_options WHERE question_id = $1 ORDER BY id ASC`,
			qID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var opt QuizOption
			if err := rows.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText, &opt.CreatedAt, &opt.UpdatedAt); err != nil {
				return nil, err
			}
			q.Options = append(q.Options, opt)
		}
	}

	return q, nil
}

func (r *QuizRepository) CreateAttempt(attempt *QuizAttempt) error {
	tx, err := r.raw.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var questionIDs []int64
	qRows, err := tx.Query(`SELECT id FROM quiz_questions WHERE quiz_id = $1 ORDER BY id ASC`, attempt.QuizID)
	if err != nil {
		return err
	}
	defer qRows.Close()

	for qRows.Next() {
		var id int64
		if err := qRows.Scan(&id); err != nil {
			return err
		}
		questionIDs = append(questionIDs, id)
	}

	source := mathrand.NewSource(time.Now().UnixNano())
	rng := mathrand.New(source)

	if len(questionIDs) > 1 {
		rng.Shuffle(len(questionIDs), func(i, j int) {
			questionIDs[i], questionIDs[j] = questionIDs[j], questionIDs[i]
		})
	}
	attempt.QuestionIDs = questionIDs

	optionOrder := make(map[string][]int64)
	for _, qID := range questionIDs {
		var optIDs []int64
		rows, err := tx.Query(`SELECT id FROM quiz_options WHERE question_id = $1`, qID)
		if err != nil {
			return err
		}
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			optIDs = append(optIDs, id)
		}
		rows.Close()

		if len(optIDs) > 1 {
			rng.Shuffle(len(optIDs), func(i, j int) {
				optIDs[i], optIDs[j] = optIDs[j], optIDs[i]
			})
		}
		optionOrder[strconv.FormatInt(qID, 10)] = optIDs
	}
	optOrderJSON, _ := json.Marshal(optionOrder)
	attempt.OptionOrder = optOrderJSON

	query := `
		INSERT INTO quiz_attempts (quiz_id, user_id, status, question_ids, option_order, current_question_index, started_at)
		VALUES ($1, $2, 'in_progress', $3, $4, 0, NOW())
		RETURNING id, started_at, created_at, updated_at
	`
	err = tx.QueryRow(query, attempt.QuizID, attempt.UserID, pq.Array(attempt.QuestionIDs), attempt.OptionOrder).Scan(
		&attempt.ID, &attempt.StartedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
	)
	if err != nil {
		return err
	}
	attempt.CurrentQuestionIndex = 0

	return tx.Commit()
}

func (r *QuizRepository) SubmitAnswers(attemptID uint, answers []QuizAnswer) (*QuizAttempt, error) {
	tx, err := r.raw.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	insertAnswerQuery := `
		INSERT INTO quiz_answers (attempt_id, question_id, option_id, is_correct)
		VALUES ($1, $2, $3, (SELECT is_correct FROM quiz_options WHERE id = $3))
		ON CONFLICT (attempt_id, question_id) DO UPDATE
			SET option_id  = EXCLUDED.option_id,
			    is_correct = EXCLUDED.is_correct
	`
	for _, ans := range answers {
		if _, err := tx.Exec(insertAnswerQuery, attemptID, ans.QuestionID, ans.OptionID); err != nil {
			return nil, err
		}
	}

	scoreQuery := `
		SELECT
			ROUND(
				(COUNT(*) FILTER (WHERE a.is_correct = TRUE)::NUMERIC
				/ NULLIF(COUNT(*), 0)) * 100,
			2) AS score
		FROM quiz_answers a
		WHERE a.attempt_id = $1
	`
	var score float64
	if err := tx.QueryRow(scoreQuery, attemptID).Scan(&score); err != nil {
		return nil, err
	}

	attempt := new(QuizAttempt)
	updateQuery := `
		UPDATE quiz_attempts
		SET status       = 'completed',
		    score        = $1,
		    submitted_at = NOW(),
		    updated_at   = NOW()
		WHERE id = $2
		RETURNING id, quiz_id, user_id, status, score, started_at, submitted_at, created_at, updated_at
	`
	if err := tx.QueryRow(updateQuery, score, attemptID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.Status, &attempt.Score,
		&attempt.StartedAt, &attempt.SubmittedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return attempt, tx.Commit()
}

func (r *QuizRepository) ResetAttempt(userID, quizID uint, adminID uint) error {
	query := `
		UPDATE quiz_attempts
		SET status     = 'reset',
		    reset_at   = NOW(),
		    reset_by   = $1,
		    updated_at = NOW()
		WHERE user_id = $2 AND quiz_id = $3 AND status != 'reset'
	`
	result, err := r.db.Exec(query, adminID, userID, quizID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *QuizRepository) GetAttemptByID(attemptID uint) (*QuizAttempt, []QuizAnswer, error) {
	attemptQuery := `
		SELECT id, quiz_id, user_id, status, score, question_ids, started_at, submitted_at, reset_at, reset_by, created_at, updated_at
		FROM quiz_attempts
		WHERE id = $1
	`
	attempt := new(QuizAttempt)
	if err := r.db.QueryRow(attemptQuery, attemptID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.Status, &attempt.Score,
		&attempt.QuestionIDs, &attempt.OptionOrder, &attempt.StartedAt, &attempt.SubmittedAt, &attempt.ResetAt,
		&attempt.ResetBy, &attempt.CreatedAt, &attempt.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	answersQuery := `
		SELECT id, attempt_id, question_id, option_id, is_correct, created_at
		FROM quiz_answers
		WHERE attempt_id = $1
	`
	rows, err := r.db.Query(answersQuery, attemptID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	answers := make([]QuizAnswer, 0)
	for rows.Next() {
		ans := QuizAnswer{}
		if err := rows.Scan(&ans.ID, &ans.AttemptID, &ans.QuestionID, &ans.OptionID, &ans.IsCorrect, &ans.CreatedAt); err != nil {
			return nil, nil, err
		}
		answers = append(answers, ans)
	}

	return attempt, answers, nil
}

func (r *QuizRepository) GetStudentQuizHistories(userID uint) ([]*QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, status, score, question_ids, option_order,
		       current_question_index, started_at, submitted_at, reset_at, reset_by,
		       created_at, updated_at
		FROM quiz_attempts
		WHERE user_id = $1 AND status IN ('completed', 'reset')
		ORDER BY COALESCE(submitted_at, reset_at) DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attempts := make([]*QuizAttempt, 0)
	for rows.Next() {
		attempt := new(QuizAttempt)
		if err := rows.Scan(
			&attempt.ID, &attempt.QuizID, &attempt.UserID, &attempt.Status, &attempt.Score,
			&attempt.QuestionIDs, &attempt.OptionOrder, &attempt.CurrentQuestionIndex,
			&attempt.StartedAt, &attempt.SubmittedAt, &attempt.ResetAt, &attempt.ResetBy,
			&attempt.CreatedAt, &attempt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	return attempts, nil
}

type QuizAdminRepository struct {
	db facades.DBExecutor
}

func NewQuizAdminRepository(db facades.DBExecutor) *QuizAdminRepository {
	return &QuizAdminRepository{db: db}
}

func (r *QuizAdminRepository) WithExecutor(executor facades.DBExecutor) *QuizAdminRepository {
	return &QuizAdminRepository{db: executor}
}

func (r *QuizAdminRepository) ListQuizzes() ([]*QuizQuiz, error) {
	query := `
		SELECT id, code, title, description, passing_score, time_limit_minutes, is_active, created_at, updated_at
		FROM quiz_quizzes
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzes := make([]*QuizQuiz, 0)
	for rows.Next() {
		q := new(QuizQuiz)
		if err := rows.Scan(
			&q.ID, &q.Code, &q.Title, &q.Description, &q.PassingScore,
			&q.TimeLimitMinutes, &q.IsActive, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, nil
}

func (r *QuizAdminRepository) GetQuizDetail(quizID uint) (*QuizQuiz, []QuizQuestion, error) {
	quizQuery := `
		SELECT id, code, title, description, passing_score, time_limit_minutes, is_active, created_at, updated_at
		FROM quiz_quizzes
		WHERE id = $1 AND deleted_at IS NULL
	`
	quiz := new(QuizQuiz)
	err := r.db.QueryRow(quizQuery, quizID).Scan(
		&quiz.ID, &quiz.Code, &quiz.Title, &quiz.Description, &quiz.PassingScore,
		&quiz.TimeLimitMinutes, &quiz.IsActive, &quiz.CreatedAt, &quiz.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	questionsQuery := `
		SELECT id, quiz_id, question_text, created_at, updated_at
		FROM quiz_questions
		WHERE quiz_id = $1
		ORDER BY id ASC
	`
	qRows, err := r.db.Query(questionsQuery, quizID)
	if err != nil {
		return nil, nil, err
	}
	defer qRows.Close()

	questions := make([]QuizQuestion, 0)
	questionIDs := make([]int64, 0)
	questionMap := make(map[uint]*QuizQuestion)

	for qRows.Next() {
		q := QuizQuestion{}
		if err := qRows.Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, nil, err
		}
		q.Options = make([]QuizOption, 0)
		questions = append(questions, q)
		questionIDs = append(questionIDs, int64(q.ID))
	}

	if len(questionIDs) > 0 {
		optionsQuery := `
			SELECT id, question_id, option_text, is_correct, created_at, updated_at
			FROM quiz_options
			WHERE question_id = ANY($1)
			ORDER BY question_id, id ASC
		`
		oRows, err := r.db.Query(optionsQuery, pq.Array(questionIDs))
		if err != nil {
			return nil, nil, err
		}
		defer oRows.Close()

		for _, q := range questions {
			qCopy := q
			questionMap[q.ID] = &qCopy
		}

		for oRows.Next() {
			opt := QuizOption{}
			if err := oRows.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText, &opt.IsCorrect, &opt.CreatedAt, &opt.UpdatedAt); err != nil {
				return nil, nil, err
			}
			if q, ok := questionMap[opt.QuestionID]; ok {
				q.Options = append(q.Options, opt)
			}
		}

		for i, q := range questions {
			if updated, ok := questionMap[q.ID]; ok {
				questions[i] = *updated
			}
		}
	}

	return quiz, questions, nil
}

func (r *QuizAdminRepository) CreateQuiz(quiz *QuizQuiz) error {
	code, err := generateQuizCode()
	if err != nil {
		return fmt.Errorf("failed to generate quiz code: %w", err)
	}
	quiz.Code = code

	query := `
		INSERT INTO quiz_quizzes (code, title, description, passing_score, time_limit_minutes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		quiz.Code, quiz.Title, quiz.Description, quiz.PassingScore, quiz.TimeLimitMinutes, quiz.IsActive,
	).Scan(&quiz.ID, &quiz.CreatedAt, &quiz.UpdatedAt)
}

func (r *QuizAdminRepository) UpdateQuiz(quiz *QuizQuiz) error {
	query := `
		UPDATE quiz_quizzes
		SET title               = $1,
		    description         = $2,
		    passing_score       = $3,
		    time_limit_minutes  = $4,
		    is_active           = $5,
		    updated_at          = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`
	result := r.db.QueryRow(query,
		quiz.Title, quiz.Description, quiz.PassingScore, quiz.TimeLimitMinutes, quiz.IsActive, quiz.ID,
	)
	return result.Scan(&quiz.UpdatedAt)
}

func (r *QuizAdminRepository) DeleteQuiz(quizID uint) error {
	query := `UPDATE quiz_quizzes SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Exec(query, quizID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *QuizAdminRepository) CreateQuestion(q *QuizQuestion) error {
	query := `
		INSERT INTO quiz_questions (quiz_id, question_text)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, q.QuizID, q.QuestionText).
		Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
}

func (r *QuizAdminRepository) UpdateQuestion(q *QuizQuestion) error {
	query := `
		UPDATE quiz_questions
		SET question_text = $1,
		    updated_at    = NOW()
		WHERE id = $2 AND quiz_id = $3
		RETURNING updated_at
	`
	return r.db.QueryRow(query, q.QuestionText, q.ID, q.QuizID).Scan(&q.UpdatedAt)
}

func (r *QuizAdminRepository) DeleteQuestion(questionID, quizID uint) error {
	query := `DELETE FROM quiz_questions WHERE id = $1 AND quiz_id = $2`
	res, err := r.db.Exec(query, questionID, quizID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *QuizAdminRepository) CreateOption(opt *QuizOption) error {
	query := `
		INSERT INTO quiz_options (question_id, option_text, is_correct)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, opt.QuestionID, opt.OptionText, opt.IsCorrect).
		Scan(&opt.ID, &opt.CreatedAt, &opt.UpdatedAt)
}

func (r *QuizAdminRepository) UpdateOption(opt *QuizOption) error {
	query := `
		UPDATE quiz_options
		SET option_text  = $1,
		    is_correct   = $2,
		    updated_at   = NOW()
		WHERE id = $3 AND question_id = $4
		RETURNING updated_at
	`
	return r.db.QueryRow(query, opt.OptionText, opt.IsCorrect, opt.ID, opt.QuestionID).
		Scan(&opt.UpdatedAt)
}

func (r *QuizAdminRepository) DeleteOption(optionID, questionID uint) error {
	query := `DELETE FROM quiz_options WHERE id = $1 AND question_id = $2`
	res, err := r.db.Exec(query, optionID, questionID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *QuizAdminRepository) ListAttempts(quizID uint) ([]QuizAttemptWithUser, error) {
	query := `
		SELECT
			a.id, a.quiz_id, a.user_id, a.status, a.score,
			a.started_at, a.submitted_at, a.reset_at, a.reset_by,
			a.created_at, a.updated_at,
			u.name, u.email
		FROM quiz_attempts a
			LEFT JOIN users u ON u.id = a.user_id
		WHERE a.quiz_id = $1
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attempts := make([]QuizAttemptWithUser, 0)
	for rows.Next() {
		var a QuizAttemptWithUser
		if err := rows.Scan(
			&a.ID, &a.QuizID, &a.UserID, &a.Status, &a.Score,
			&a.StartedAt, &a.SubmittedAt, &a.ResetAt, &a.ResetBy,
			&a.CreatedAt, &a.UpdatedAt,
			&a.UserName, &a.UserEmail,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

type QuizAttemptWithUser struct {
	QuizAttempt
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
