package models

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type MeetingSession struct {
	ID                     int       `json:"id"`
	UserID                 int       `json:"user_id"`
	MentorID               int       `json:"mentor_id"`
	User                   User      `json:"student"`
	Mentor                 User      `json:"mentor"`
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
	IsStudentAttended      bool      `json:"is_student_attended"`
	IsMentorAttended       bool      `json:"is_mentor_attended"`
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
        SELECT 
            ms.id, ms.user_id, ms.mentor_id, ms.session_date, ms.session_time, 
            ms.session_duration, ms.session_type, ms.session_topic, ms.session_description,
            ms.session_proof, ms.session_feedback, ms.student_attendance_proof,
            ms.mentor_attendance_proof, ms.session_status, ms.is_student_attended, ms.is_mentor_attended, ms.created_at, ms.updated_at,
            u.id as user_id, u.name as user_name, u.email as user_email, u.role as user_role,
            m.id as mentor_id, m.name as mentor_name, m.email as mentor_email, m.role as mentor_role
        FROM meeting_sessions ms
        LEFT JOIN users u ON ms.user_id = u.id
        LEFT JOIN users m ON ms.mentor_id = m.id
        ORDER BY ms.created_at DESC
    `

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*MeetingSession, 0)
	for rows.Next() {
		session := &MeetingSession{}
		user := &User{}
		mentor := &User{}

		err := rows.Scan(
			&session.ID, &session.UserID, &session.MentorID, &session.SessionDate,
			&session.SessionTime, &session.SessionDuration, &session.SessionType,
			&session.SessionTopic, &session.SessionDescription, &session.SessionProof,
			&session.SessionFeedback, &session.StudentAttendanceProof,
			&session.MentorAttendanceProof, &session.SessionStatus,
			&session.IsStudentAttended, &session.IsMentorAttended,
			&session.CreatedAt, &session.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.Role,
			&mentor.ID, &mentor.Name, &mentor.Email, &mentor.Role,
		)
		if err != nil {
			return nil, err
		}

		session.User = *user
		session.Mentor = *mentor
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (repo *MeetingSessionRepository) GetByID(id int) (*MeetingSession, error) {
	query := `
        SELECT 
            ms.id, ms.user_id, ms.mentor_id, ms.session_date, ms.session_time, 
            ms.session_duration, ms.session_type, ms.session_topic, ms.session_description,
            ms.session_proof, ms.session_feedback, ms.student_attendance_proof,
            ms.mentor_attendance_proof, ms.session_status, ms.is_student_attended, ms.is_mentor_attended, ms.created_at, ms.updated_at,
            u.id, u.name, u.email, u.role, u.created_at, u.updated_at,
            m.id, m.name, m.email, m.role, m.created_at, m.updated_at
        FROM meeting_sessions ms
        LEFT JOIN users u ON ms.user_id = u.id
        LEFT JOIN users m ON ms.mentor_id = m.id
        WHERE ms.id = $1
    `

	session := &MeetingSession{}
	user := &User{}
	mentor := &User{}

	err := repo.db.QueryRow(query, id).Scan(
		&session.ID, &session.UserID, &session.MentorID, &session.SessionDate,
		&session.SessionTime, &session.SessionDuration, &session.SessionType,
		&session.SessionTopic, &session.SessionDescription, &session.SessionProof,
		&session.SessionFeedback, &session.StudentAttendanceProof,
		&session.MentorAttendanceProof, &session.SessionStatus,
		&session.IsStudentAttended, &session.IsMentorAttended,
		&session.CreatedAt, &session.UpdatedAt,
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		&mentor.ID, &mentor.Name, &mentor.Email, &mentor.Role, &mentor.CreatedAt, &mentor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	session.User = *user
	session.Mentor = *mentor

	return session, nil
}

func (repo *MeetingSessionRepository) GetByUser(userID int) ([]*MeetingSession, error) {
	query := `
		SELECT 
			ms.id, ms.user_id, ms.mentor_id, ms.session_date, ms.session_time, 
			ms.session_duration, ms.session_type, ms.session_topic, ms.session_description,
			ms.session_proof, ms.session_feedback, ms.student_attendance_proof,
			ms.mentor_attendance_proof, ms.session_status, ms.is_student_attended, ms.is_mentor_attended, ms.created_at, ms.updated_at,
			u.id as user_id, u.name as user_name, u.email as user_email, u.role as user_role,
			m.id as mentor_id, m.name as mentor_name, m.email as mentor_email, m.role as mentor_role
		FROM meeting_sessions ms
		LEFT JOIN users u ON ms.user_id = u.id
		LEFT JOIN users m ON ms.mentor_id = m.id
		WHERE (u.id = $1 OR m.id = $1)
		ORDER BY ms.created_at DESC
	`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*MeetingSession, 0)
	for rows.Next() {
		session := &MeetingSession{}
		user := &User{}
		mentor := &User{}

		err := rows.Scan(
			&session.ID, &session.UserID, &session.MentorID, &session.SessionDate,
			&session.SessionTime, &session.SessionDuration, &session.SessionType,
			&session.SessionTopic, &session.SessionDescription, &session.SessionProof,
			&session.SessionFeedback, &session.StudentAttendanceProof,
			&session.MentorAttendanceProof, &session.SessionStatus,
			&session.IsStudentAttended, &session.IsMentorAttended,
			&session.CreatedAt, &session.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.Role,
			&mentor.ID, &mentor.Name, &mentor.Email, &mentor.Role,
		)
		if err != nil {
			return nil, err
		}

		session.User = *user
		session.Mentor = *mentor
		sessions = append(sessions, session)
	}

	return sessions, nil
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

func (repo *MeetingSessionRepository) UpdateProofs(id int, sessionProof, studentAttendanceProof, mentorAttendanceProof, sessionFeedback *string) error {
	setClauses := []string{}
	args := []interface{}{}

	if sessionProof != nil {
		setClauses = append(setClauses, "session_proof = ?")
		args = append(args, *sessionProof)
	}
	if studentAttendanceProof != nil {
		setClauses = append(setClauses, "student_attendance_proof = ?")
		args = append(args, *studentAttendanceProof)
	}
	if sessionProof != nil || studentAttendanceProof != nil {
		setClauses = append(setClauses, "is_student_attended = TRUE")
	}
	if mentorAttendanceProof != nil {
		setClauses = append(setClauses, "mentor_attendance_proof = ?")
		args = append(args, *mentorAttendanceProof)
	}
	if sessionFeedback != nil {
		setClauses = append(setClauses, "session_feedback = ?")
		args = append(args, *sessionFeedback)
	}
	if sessionFeedback != nil || mentorAttendanceProof != nil {
		setClauses = append(setClauses, "is_mentor_attended = TRUE")
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := "UPDATE meeting_sessions SET " +
		strings.Join(setClauses, ", ") +
		" WHERE id = ?"
	args = append(args, id)

	for i := range args {
		query = strings.Replace(query, "?", "$"+strconv.Itoa(i+1), 1)
	}

	_, err := repo.db.Exec(query, args...)
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

func (repo *MeetingSessionRepository) VerifyAttendance(id int, userId int, isMentor bool) error {
	var column string
	if isMentor {
		column = "mentor_attendance_proof"
	} else {
		column = "student_attendance_proof"
	}

	var idColumn string
	if isMentor {
		idColumn = "mentor_id"
	} else {
		idColumn = "user_id"
	}

	query := `
		SELECT id, student_attendance_proof, mentor_attendance_proof FROM meeting_sessions
		WHERE id = $1 AND ` + idColumn + ` = $2 AND ` + column + ` IS NULL
		AND session_status = 'scheduled'
		AND session_date >= CURRENT_DATE
		AND (
			session_date > CURRENT_DATE OR
			(session_date = CURRENT_DATE AND session_time >= CURRENT_TIME)
		)
	`

	meetingSession := new(MeetingSession)
	err := repo.db.QueryRow(query, id, userId).Scan(&meetingSession.ID, &meetingSession.StudentAttendanceProof, &meetingSession.MentorAttendanceProof)

	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}

		return err
	}

	if meetingSession.ID == 0 {
		return sql.ErrNoRows
	}

	if !isMentor && meetingSession.StudentAttendanceProof != nil {
		return sql.ErrNoRows
	}

	if isMentor && meetingSession.MentorAttendanceProof != nil {
		return sql.ErrNoRows
	}

	return nil
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
