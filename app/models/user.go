package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
	gomail "github.com/studio-senkou/lentera-cendekia-be/utils/mail"
	"github.com/studio-senkou/lentera-cendekia-be/utils/queue"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"` // 'user', 'mentor', 'admin'
	Password        string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
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

func (u *User) MarkAsActive() {
	u.IsActive = true
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
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.EmailVerifiedAt, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) GetUserCount() (map[string]int, error) {
	query := `SELECT role, COUNT(*) FROM users WHERE is_active = true GROUP BY role`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println(rows)

	userCount := make(map[string]int)
	for rows.Next() {
		var role string
		var count int

		if err := rows.Scan(&role, &count); err != nil {
			return nil, err
		}

		if role != "admin" {
			if role == "user" {
				userCount["student"] = count
			} else {
				userCount[role] = count
			}
		}
	}

	return userCount, nil
}

func (r *UserRepository) GetUserDropdown() ([]*Student, error) {
	query := `
		SELECT
			s.id, u.name, u.email
		FROM students s
			LEFT OUTER JOIN users u ON u.id = s.user_id
			WHERE u.role = 'user'
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]*Student, 0)

	for rows.Next() {

		student := new(Student)
		user := new(User)

		err := rows.Scan(&student.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		student.User = *user
		students = append(students, student)
	}

	return students, nil
}

func (r *UserRepository) GetMentorDropdown() ([]*Mentor, error) {
	query := `
		SELECT
			id, name, email
		FROM users WHERE role = 'mentor'
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentors := make([]*Mentor, 0)
	for rows.Next() {

		mentor := new(Mentor)
		user := new(User)

		err := rows.Scan(&mentor.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		mentor.User = *user
		mentors = append(mentors, mentor)
	}

	return mentors, nil
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	query := `SELECT id, name, email, role, email_verified_at, is_active, created_at, updated_at FROM users WHERE id = $1`

	user := new(User)
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.EmailVerifiedAt, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password, role, email_verified_at, is_active, created_at, updated_at FROM users WHERE email = $1`

	user := new(User)
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.EmailVerifiedAt, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(user *User) error {
	query := `INSERT INTO users (name, email, password, role, email_verified_at, is_active) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, created_at, updated_at`

	hashedPassword, err := user.HashPassword()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		hashedPassword,
		user.Role,
		user.EmailVerifiedAt,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {

		if strings.Contains(err.Error(), "users_email_key") {
			return ErrEmailAlreadyExists
		}

		return errors.New("failed to create user: " + err.Error())
	}

	if user.EmailVerifiedAt == nil && !user.IsActive {
		activationToken, err := auth.GenerateOneTimeToken(user.ID, "account_activation", 24*time.Hour)
		if err != nil {
			return errors.New("failed to generate activation token")
		}

		client := queue.NewClient()

		client.RegisterHandlerFunc("email:send", func(ctx context.Context, task *asynq.Task) error {

			email, err := gomail.NewMailFromTemplate(
				user.Email,
				"Welcome aboard to Lentera Cendekia",
				"templates/emails/welcome.html",
				fiber.Map{
					"Name":           user.Name,
					"ActivationLink": fmt.Sprintf("%s/activate?token=%s", app.GetEnv("APP_FE_URL", "http://localhost:3000"), activationToken.Token),
				},
			)
			if err != nil {
				return errors.New("failed to create welcome email")
			}

			return email.Send()

		})
	}

	return err
}

func (r *UserRepository) Update(user *User) (string, error) {
	current, err := r.GetByID(user.ID)
	if err != nil {
		return "", err
	}

	var (
		setClauses   []string
		args         []any
		argIdx       = 1
		emailChanged bool
	)

	if user.Name != "" && user.Name != current.Name {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argIdx))
		args = append(args, user.Name)
		argIdx++
	}
	if user.Email != "" && user.Email != current.Email {
		setClauses = append(setClauses, "email = $"+strconv.Itoa(argIdx))
		args = append(args, user.Email)
		argIdx++
		emailChanged = true
	}
	if user.EmailVerifiedAt != nil && current.EmailVerifiedAt == nil {
		setClauses = append(setClauses, "email_verified_at = $"+strconv.Itoa(argIdx))
		args = append(args, user.EmailVerifiedAt)
		argIdx++
	}
	if user.IsActive != current.IsActive {
		setClauses = append(setClauses, "is_active = $"+strconv.Itoa(argIdx))
		args = append(args, user.IsActive)
		argIdx++
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	if len(setClauses) == 1 {
		return "", nil
	}

	setClause := strings.Join(setClauses, ", ")
	query := `UPDATE users SET ` + setClause + ` WHERE id = $` + strconv.Itoa(argIdx) + ` RETURNING updated_at`
	args = append(args, user.ID)

	fmt.Println("Executing query:", query, "with args:", args)

	err = r.db.QueryRow(query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		return "", err
	}

	if emailChanged {
		return user.Email, nil
	}
	return "", nil
}

func (r *UserRepository) UpdatePassword(id uint, newPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, hashedPassword, id)
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

func (r *UserRepository) VerifyOldPassword(id int, oldPassword string) (bool, error) {
	query := `SELECT password FROM users WHERE id = $1`

	var hashedPassword string
	err := r.db.QueryRow(query, id).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		return false, nil
	}

	return true, nil
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
