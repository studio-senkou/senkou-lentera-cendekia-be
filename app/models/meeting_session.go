package models

import (
	"database/sql"
	"time"
)

type MeetingSession struct {
	ID          uint       `json:"id"`
	StudentID   uint       `json:"student_id"`
	Student     Student    `json:"student"`
	MentorID    uint       `json:"mentor_id"`
	Mentor      Mentor     `json:"mentor"`
	Date        DateOnly   `json:"session_date"`
	Time        TimeOnly   `json:"session_time"`
	Duration    uint       `json:"duration_minutes"`
	Status      string     `json:"status"` // "pending", "confirmed", "completed", "canceled"
	Description string     `json:"description"`
	Note        *string    `json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type MeetingSessionRepository struct {
	db *sql.DB
}

func NewMeetingSessionRepository(db *sql.DB) *MeetingSessionRepository {
	return &MeetingSessionRepository{
		db: db,
	}
}

func (r *MeetingSessionRepository) Create(session *MeetingSession) (*MeetingSession, error) {

	query := `
		INSERT INTO meeting_sessions (
			student_id,
			mentor_id,
			session_date,
			session_time,
			duration_minutes,
			status,
			note,
			description
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at, updated_at
	`

	if err := r.db.QueryRow(query,
		session.StudentID,
		session.MentorID,
		session.Date,
		session.Time,
		session.Duration,
		session.Status,
		session.Note,
		session.Description,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *MeetingSessionRepository) BulkCreateSessions(sessions []*MeetingSession) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, session := range sessions {
		query := `
			INSERT INTO meeting_sessions (
				student_id,
				mentor_id,
				session_date,
				session_time,
				duration_minutes,
				status,
				note,
				description
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8
			) RETURNING id, created_at, updated_at
		`

		if err := tx.QueryRow(query,
			session.StudentID,
			session.MentorID,
			session.Date,
			session.Time,
			session.Duration,
			session.Status,
			session.Note,
			session.Description,
		).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MeetingSessionRepository) GetAll(userID uint) ([]*MeetingSession, error) {

	query := `
		SELECT 
			ms.id, ms.student_id, ms.mentor_id, ms.session_date, ms.session_time,
			ms.duration_minutes, ms.status, ms.note, ms.description, ms.created_at, ms.updated_at, ms.deleted_at,
			u.id, u.name, u.email, mu.id, mu.name, mu.email
		FROM meeting_sessions ms
			LEFT JOIN students s ON s.id = ms.student_id
			LEFT JOIN users u ON u.id = s.user_id
			LEFT JOIN mentors m ON m.id = ms.mentor_id
			LEFT JOIN users mu ON mu.id = m.user_id
	`

	var rows *sql.Rows
	var err error

	if userID != 0 {
		query += " WHERE s.user_id = $1"
		rows, err = r.db.Query(query, userID)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetingSessions := make([]*MeetingSession, 0)

	for rows.Next() {
		session := new(MeetingSession)
		student := new(Student)
		mentor := new(Mentor)

		if err := rows.Scan(
			&session.ID, &session.StudentID, &session.MentorID, &session.Date, &session.Time,
			&session.Duration, &session.Status, &session.Note, &session.Description, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
			&student.ID, &student.User.Name, &student.User.Email,
			&mentor.ID, &mentor.User.Name, &mentor.User.Email,
		); err != nil {
			return nil, err
		}

		session.Student = *student
		session.Mentor = *mentor

		meetingSessions = append(meetingSessions, session)
	}

	return meetingSessions, nil
}

func (r *MeetingSessionRepository) GetByID(id uint) (*MeetingSession, error) {
	query := `
		SELECT 
			ms.id, ms.student_id, ms.mentor_id, ms.session_date, ms.session_time,
			ms.duration_minutes, ms.status, ms.note, ms.description, ms.created_at, ms.updated_at, ms.deleted_at,
			u.id, u.name, u.email, mu.id, mu.name, mu.email
		FROM meeting_sessions ms
			LEFT JOIN students s ON s.id = ms.student_id
			LEFT JOIN users u ON u.id = s.user_id
			LEFT JOIN mentors m ON m.id = ms.mentor_id
			LEFT JOIN users mu ON mu.id = m.user_id
		WHERE ms.id = $1
	`

	row := r.db.QueryRow(query, id)

	session := new(MeetingSession)
	student := new(Student)
	mentor := new(Mentor)

	if err := row.Scan(
		&session.ID, &session.StudentID, &session.MentorID, &session.Date, &session.Time,
		&session.Duration, &session.Status, &session.Note, &session.Description, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
		&student.ID, &student.User.Name, &student.User.Email,
		&mentor.ID, &mentor.User.Name, &mentor.User.Email,
	); err != nil {
		return nil, err
	}

	session.Student = *student
	session.Mentor = *mentor

	return session, nil
}

func (r *MeetingSessionRepository) BulkUpdate(sessions []*MeetingSession) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, session := range sessions {
		query := `
			UPDATE meeting_sessions SET
				student_id = $1,
				mentor_id = $2,
				session_date = $3,
				session_time = $4,
				duration_minutes = $5,
				status = $6,
				note = $7,
				description = $8,
				updated_at = NOW()
			WHERE id = $9
		`

		_, err := tx.Exec(query,
			session.StudentID,
			session.MentorID,
			session.Date,
			session.Time,
			session.Duration,
			session.Status,
			session.Note,
			session.Description,
			session.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MeetingSessionRepository) Delete(id uint) error {
	query := `
		DELETE FROM meeting_sessions WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}
