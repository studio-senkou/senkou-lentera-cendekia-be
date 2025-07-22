package models

import (
	"database/sql"
)

type UserHasToken struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

type AuthenticationRepository struct {
	db *sql.DB
}

func NewAuthenticationRepository(db *sql.DB) *AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func (r *AuthenticationRepository) Create(userHasToken *UserHasToken) error {
	query := `INSERT INTO user_has_tokens (user_id, token) VALUES ($1, $2)`
	_, err := r.db.Exec(query, userHasToken.UserID, userHasToken.Token)
	return err
}

func (r *AuthenticationRepository) UpdateOrCreate(userHasToken *UserHasToken) error {
	query := `INSERT INTO user_has_tokens (user_id, token) 
			  VALUES ($1, $2) 
			  ON CONFLICT (user_id) 
			  DO UPDATE SET token = EXCLUDED.token, updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, userHasToken.UserID, userHasToken.Token)
	return err
}

func (r *AuthenticationRepository) ValidateSessionExist(userID int, token string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_has_tokens WHERE user_id = $1 AND token = $2)`

	var exists bool
	err := r.db.QueryRow(query, userID, token).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthenticationRepository) InvalidateToken(userID int) error {
	query := `DELETE FROM user_has_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
