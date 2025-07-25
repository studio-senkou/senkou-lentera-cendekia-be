package models

import (
	"database/sql"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"` // 'user', 'mentor', 'admin'
	Password        string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (u *User) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) MarkEmailAsVerified() {
	now := time.Now()
	u.EmailVerifiedAt = &now
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]*User, error) {
	query := `
		SELECT  id, name, email, role, email_verified_at,  is_active, created_at, updated_at 
		FROM users 
		WHERE role NOT IN ('admin')
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) GetUserDropdown() ([]*User, error) {
	query := `SELECT id, name FROM users WHERE role = 'user' ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) GetMentorDropdown() ([]*User, error) {
	query := `SELECT id, name FROM users WHERE role = 'mentor' ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	user := new(User)
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = $1`

	user := new(User)
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(user *User) error {
	query := `INSERT INTO users (name, email, password, role, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, NOW(), NOW()) 
			  RETURNING id, created_at, updated_at`

	hashedPassword, err := user.HashPassword()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, user.Name, user.Email, hashedPassword, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *UserRepository) Update(user *User) (string, error) {
	updatedUser, err := r.GetByID(user.ID)
	if err != nil {
		return "", err
	}

	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if user.Name != "" && user.Name != updatedUser.Name {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argIdx))
		args = append(args, user.Name)
		argIdx++
	}

	emailChanged := false
	if user.Email != "" && user.Email != updatedUser.Email {
		setClauses = append(setClauses, "email = $"+strconv.Itoa(argIdx))
		args = append(args, user.Email)
		argIdx++
		emailChanged = true
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	if len(setClauses) == 1 {
		return "", nil
	}

	setClause := ""
	for i, clause := range setClauses {
		if i > 0 {
			setClause += ", "
		}
		setClause += clause
	}

	query := `UPDATE users SET ` + setClause + ` WHERE id = $` + strconv.Itoa(argIdx) + ` RETURNING updated_at`
	args = append(args, user.ID)

	err = r.db.QueryRow(query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		return "", err
	}

	if emailChanged {
		return user.Email, nil
	}
	return "", nil
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

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
