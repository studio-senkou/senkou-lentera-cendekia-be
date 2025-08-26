package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID        uuid.UUID  `json:"ID"`
	ClassName string     `json:"classname"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{
		db: db,
	}
}

func (r *ClassRepository) Store(className string) error {
	classId := uuid.New()

	query := `
		INSERT INTO classes (id, classname) 
		VALUES ($1, $2) RETURNING id
	`

	err := r.db.QueryRow(query, classId, className).Err()

	return err
}

func (r *ClassRepository) FindAll() ([]*Class, error) {
	query := `
		SELECT 
			id,
			classname,
			created_at
		FROM classes	
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes := make([]*Class, 0)
	for rows.Next() {
		var class Class

		err := rows.Scan(&class.ID, &class.ClassName, &class.CreatedAt)

		if err != nil {
			return nil, err
		}

		classes = append(classes, &class)
	}

	return classes, nil
}
