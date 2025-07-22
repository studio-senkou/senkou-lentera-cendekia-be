package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // 'user', 'mentor', 'admin'
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]*User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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
	query := `INSERT INTO users (name, email, password, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW()) 
			  RETURNING id, created_at, updated_at`

	hashedPassword, err := user.HashPassword()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, user.Name, user.Email, hashedPassword).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *UserRepository) Update(user *User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() 
			  WHERE id = $3 
			  RETURNING updated_at`

	err := r.db.QueryRow(query, user.Name, user.Email, user.ID).Scan(&user.UpdatedAt)
	return err
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
