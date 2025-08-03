package models

import "database/sql"

type Blog struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	AuthorID  int     `json:"author_id"`
	Author    User    `json:"author"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *Blog) error {
	query := `
		INSERT INTO blogs (title, content, author_id, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, blog.Title, blog.Content, blog.AuthorID).Scan(&blog.ID, &blog.CreatedAt, &blog.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *BlogRepository) GetByID(id int) (*Blog, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at, deleted_at
		FROM blogs
		WHERE id = ? AND deleted_at IS NULL`

	blog := &Blog{}
	err := r.db.QueryRow(query, id).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.AuthorID, &blog.CreatedAt, &blog.UpdatedAt, &blog.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return blog, nil
}

func (r *BlogRepository) GetAll() ([]*Blog, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at, deleted_at
		FROM blogs
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		blog := &Blog{}
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.AuthorID, &blog.CreatedAt, &blog.UpdatedAt, &blog.DeletedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *BlogRepository) Update(blog *Blog) error {
	query := `
		UPDATE blogs
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.Exec(query, blog.Title, blog.Content, blog.ID)
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

func (r *BlogRepository) Delete(id int) error {
	query := `
		UPDATE blogs
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`

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

func (r *BlogRepository) Restore(id int) error {
	query := `
		UPDATE blogs
		SET deleted_at = NULL
		WHERE id = ? AND deleted_at IS NOT NULL`

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
