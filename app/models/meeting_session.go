package models

import (
	"database/sql"
	"time"
)

type MeetingSession struct {
	ID                     int       `json:"id"`
	UserID                 int       `json:"user_id"`
	MentorID               int       `json:"mentor_id"`
	User                   User      `json:"-"`
	Mentor                 User      `json:"-"`
	SessionDate            string    `json:"session_date"`
	SessionTime            string    `json:"session_time"`
	SessionDuration        int       `json:"session_duration"` // Duration in minutes
	SessionType            string    `json:"session_type"`
	SessionTopic           string    `json:"session_topic"`
	SessionDescription     *string   `json:"session_description"`
	SessionProof           *string   `json:"session_proof"`
	SessionFeedback        *string   `json:"session_feedback"`
	StudentAttendanceProof *string   `json:"student_attendance_proof"`
	MentorAttendanceProof  *string   `json:"mentor_attendance_proof"`
	SessionStatus          string    `json:"session_status"` // scheduled, completed, cancelled
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type MeetingSessionRepository struct {
	db *sql.DB
}

func NewMeetingSessionRepository(db *sql.DB) *MeetingSessionRepository {
	return &MeetingSessionRepository{db: db}
}

func (repo *MeetingSessionRepository) Create(session *MeetingSession) error {
	query := `
		INSERT INTO meeting_sessions (
			user_id,
			mentor_id,
			session_date,
			session_time,
			session_duration,
			session_type,
			session_topic
		) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, session_status
	`
	err := repo.db.QueryRow(
		query,
		session.UserID,
		session.MentorID,
		session.SessionDate,
		session.SessionTime,
		session.SessionDuration,
		session.SessionType,
		session.SessionTopic,
	).Scan(&session.ID, &session.SessionStatus)

	return err
}

func (repo *MeetingSessionRepository) GetAll() ([]*MeetingSession, error) {
	query := `
		SELECT id, user_id, mentor_id, session_date, session_time, session_duration,
		       session_type, session_topic, session_description, session_proof,
		       session_feedback, student_attendance_proof, mentor_attendance_proof,
		       session_status, created_at, updated_at
		FROM meeting_sessions
	`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*MeetingSession, 0)
	for rows.Next() {
		session := &MeetingSession{}
		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.MentorID,
			&session.SessionDate,
			&session.SessionTime,
			&session.SessionDuration,
			&session.SessionType,
			&session.SessionTopic,
			&session.SessionDescription,
			&session.SessionProof,
			&session.SessionFeedback,
			&session.StudentAttendanceProof,
			&session.MentorAttendanceProof,
			&session.SessionStatus,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (repo *MeetingSessionRepository) GetByID(id int) (*MeetingSession, error) {
	query := `
		SELECT id, user_id, mentor_id, session_date, session_time, session_duration,
		       session_type, session_topic, session_description, session_proof,
		       session_feedback, student_attendance_proof, mentor_attendance_proof,
		       session_status, created_at, updated_at
		FROM meeting_sessions WHERE id = $1
	`

	session := &MeetingSession{}
	err := repo.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.MentorID,
		&session.SessionDate,
		&session.SessionTime,
		&session.SessionDuration,
		&session.SessionType,
		&session.SessionTopic,
		&session.SessionDescription,
		&session.SessionProof,
		&session.SessionFeedback,
		&session.StudentAttendanceProof,
		&session.MentorAttendanceProof,
		&session.SessionStatus,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (repo *MeetingSessionRepository) Update(session *MeetingSession) error {
	query := `
		UPDATE meeting_sessions SET
			session_date = $1,
			session_time = $2,
			session_duration = $3,
			session_type = $4,
			session_topic = $5,
			session_description = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := repo.db.Exec(query,
		session.SessionDate,
		session.SessionTime,
		session.SessionDuration,
		session.SessionType,
		session.SessionTopic,
		session.SessionDescription,
		time.Now(),
		session.ID,
	)

	return err
}

func (repo *MeetingSessionRepository) UpdateStatus(id int, status string) error {
	var updatedStatus string
	switch status {
	case "cancel":
		updatedStatus = "cancelled"
	case "complete":
		updatedStatus = "completed"
	default:
		updatedStatus = "scheduled"
	}

	_, err := repo.db.Exec(
		`UPDATE meeting_sessions SET session_status = $1, updated_at = NOW() WHERE id = $2`,
		updatedStatus, id,
	)
	return err
}

func (repo *MeetingSessionRepository) Delete(id int) error {
	query := `DELETE FROM meeting_sessions WHERE id = $1`
	result, err := repo.db.Exec(query, id)
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
