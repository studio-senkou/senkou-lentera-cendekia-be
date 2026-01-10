package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	User      User       `json:"user"`
	ClassID   uuid.UUID  `json:"class_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{
		db: db,
	}
}

func (r *StudentRepository) AddIntoClass(userID uint, classID uuid.UUID) (*Student, error) {
	query := `
		INSERT INTO students (
			user_id,
			class_id
		) VALUES (
			$1, $2 
		) RETURNING id, created_at, updated_at
	`

	student := new(Student)

	if err := r.db.QueryRow(query, userID, classID).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) FindStudentClass(userID int) ([]*Class, error) {
	query := `
		SELECT
			c.id,
			c.classname
		FROM classes c
			INNER JOIN students s ON s.class_id = c.id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL AND c.deleted_at IS NULL
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

func (r *StudentRepository) IsInClass(userID int, classID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM students WHERE user_id = $1 AND class_id = $2 AND deleted_at IS NULL)"
	err := r.db.QueryRow(query, userID, classID).Scan(&exists)
	return exists, err
}

func (r *StudentRepository) RemoveFromClass(userID int, classID uuid.UUID) error {
	query := `UPDATE students SET deleted_at = NOW() WHERE user_id = $1 AND class_id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, userID, classID)
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
