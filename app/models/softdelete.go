package models

import (
	"database/sql"
	"fmt"
)

type SoftDeleteResult struct {
	RowsAffected int64
}

func SoftDelete(db *sql.DB, table string, idColumn string, id any) (*SoftDeleteResult, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET deleted_at = NOW() WHERE %s = $1 AND deleted_at IS NULL",
		table, idColumn,
	)

	result, err := db.Exec(query, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &SoftDeleteResult{RowsAffected: rowsAffected}, nil
}

func Restore(db *sql.DB, table string, idColumn string, id any) (*SoftDeleteResult, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET deleted_at = NULL WHERE %s = $1 AND deleted_at IS NOT NULL",
		table, idColumn,
	)

	result, err := db.Exec(query, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &SoftDeleteResult{RowsAffected: rowsAffected}, nil
}

func HardDelete(db *sql.DB, table string, idColumn string, id any) (*SoftDeleteResult, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, idColumn)

	result, err := db.Exec(query, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &SoftDeleteResult{RowsAffected: rowsAffected}, nil
}
