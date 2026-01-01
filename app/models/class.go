package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID        uuid.UUID  `json:"id"`
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

func (r *ClassRepository) Store(className string) (*Class, error) {
	classId := uuid.New()

	query := `
		INSERT INTO classes (id, classname) 
		VALUES ($1, $2) RETURNING id, classname, created_at
	`

	class := new(Class)

	if err := r.db.QueryRow(
		query,
		classId,
		className,
	).Scan(
		&class.ID,
		&class.ClassName,
		&class.CreatedAt,
	); err != nil {
		return nil, err
	}

	return class, nil
}

func (r *ClassRepository) FindAll() ([]*Class, error) {
	query := `
		SELECT 
			id,
			classname,
			created_at
		FROM classes
		WHERE deleted_at IS NULL	
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

func (r *ClassRepository) FindAllForDropdown() ([]*Class, error) {
	query := `
		SELECT 
			id,
			classname
		FROM classes
		WHERE deleted_at IS NULL
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes := make([]*Class, 0)
	for rows.Next() {
		var class Class

		err := rows.Scan(&class.ID, &class.ClassName)

		if err != nil {
			return nil, err
		}

		classes = append(classes, &class)
	}

	return classes, nil
}

func (r *ClassRepository) Update(id uuid.UUID, className string) (*Class, error) {
	query := `
		UPDATE classes
		SET classname = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, classname, created_at
	`

	class := new(Class)
	if err := r.db.QueryRow(query, className, id).Scan(&class.ID, &class.ClassName, &class.CreatedAt); err != nil {
		return nil, err
	}

	return class, nil
}

func (r *ClassRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE classes
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, id)
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