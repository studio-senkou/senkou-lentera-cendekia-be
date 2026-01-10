package models

import (
	"database/sql"
	"time"
)

type StudentPlan struct {
	ID            uint       `json:"id"`
	StudentID     uint       `json:"student_id"`
	TotalSessions uint       `json:"total_sessions"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"-"`
}

type StudentPlanRepository struct {
	db *sql.DB
}

func NewStudentPlanRepository(db *sql.DB) *StudentPlanRepository {
	return &StudentPlanRepository{
		db: db,
	}
}

func (r *StudentPlanRepository) CreateNewStudentPlan(plan *StudentPlan) error {
	query := `
		INSERT INTO student_plans (
			student_id, total_sessions
		) VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, plan.StudentID, plan.TotalSessions)
	return err
}

func (r *StudentPlanRepository) GetCurrentStudentPlan(userID uint) (*StudentPlan, error) {
	query := `
		SELECT 
			ps.id, ps.student_id, ps.total_sessions, ps.created_at, ps.updated_at, ps.deleted_at
		FROM student_plans ps
			LEFT OUTER JOIN students s ON ps.student_id = s.id
			LEFT OUTER JOIN users u ON s.user_id = u.id
		WHERE u.id = $1 AND ps.deleted_at IS NULL
		LIMIT 1
	`
	var plan StudentPlan
	err := r.db.QueryRow(query, userID).Scan(
		&plan.ID,
		&plan.StudentID,
		&plan.TotalSessions,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		&plan.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

func (r *StudentPlanRepository) UpdateStudentPlan(plan *StudentPlan) error {
	query := `
		UPDATE student_plans
		SET total_sessions = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(query, plan.TotalSessions, plan.UpdatedAt, plan.ID)
	return err
}

func (r *StudentPlanRepository) DeleteStudentPlan(planID int) error {
	query := `
		UPDATE student_plans
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, planID)
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

func (r *StudentPlanRepository) RestoreStudentPlan(planID int) error {
	query := `
		UPDATE student_plans
		SET deleted_at = NULL
		WHERE id = $1 AND deleted_at IS NOT NULL
	`
	result, err := r.db.Exec(query, planID)
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
