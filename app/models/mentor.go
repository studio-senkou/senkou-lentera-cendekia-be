package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Mentor struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	User      User       `json:"user"`
	ClassID   uuid.UUID  `json:"class_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type MentorRepository struct {
	db *sql.DB
}

func NewMentorRepository(db *sql.DB) *MentorRepository {
	return &MentorRepository{
		db: db,
	}
}

func (r *MentorRepository) AddIntoClass(userID uint, classID uuid.UUID) (*Mentor, error) {
	query := `
		INSERT INTO mentors (
			user_id,
			class_id
		) VALUES (
			$1, $2 
		) RETURNING id, created_at, updated_at
	`

	mentor := new(Mentor)

	if err := r.db.QueryRow(query, userID, classID).Scan(
		&mentor.ID,
		&mentor.CreatedAt,
		&mentor.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return mentor, nil
}

func (r *MentorRepository) FindMentorClass(userID int) ([]*Class, error) {
	query := `
		SELECT
			c.id,
			c.classname
		FROM classes c
			INNER JOIN Mentors s ON s.class_id = c.id
			WHERE s.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*Class
	for rows.Next() {
		class := &Class{}
		if err := rows.Scan(&class.ID, &class.ClassName); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (r *MentorRepository) IsInClass(userID int, classID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Mentors WHERE user_id = $1 AND class_id = $2)"
	err := r.db.QueryRow(query, userID, classID).Scan(&exists)
	return exists, err
}

func (r *MentorRepository) RemoveFromClass(userID int, classID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM Mentors WHERE user_id = $1 AND class_id = $2", userID, classID)
	return err
}
