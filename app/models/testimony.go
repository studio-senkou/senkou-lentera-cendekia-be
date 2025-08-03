package models

import "database/sql"

type Testimony struct {
	ID                         int     `json:"id"`
	TestimonerName             string  `json:"testimoner_name"`
	TestimonerCurrentPosition  string  `json:"testimoner_current_position"`
	TestimonerPreviousPosition string  `json:"testimoner_previous_position"`
	TestimonerPhoto            *string `json:"testimoner_photo"`
	TestimonyText              string  `json:"testimony_text"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  string  `json:"updated_at"`
	DeletedAt                  *string `json:"-,omitempty"`
}

type TestimonyRepository struct {
	db *sql.DB
}

func NewTestimonyRepository(db *sql.DB) *TestimonyRepository {
	return &TestimonyRepository{db: db}
}

func (r *TestimonyRepository) Create(testimony *Testimony) error {
	query := `
		INSERT INTO testimonials (
			testimoner_name,
			testimoner_current_position,
			testimoner_previous_position,
			testimoner_photo,
			testimony_text,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, testimony.TestimonerName, testimony.TestimonerCurrentPosition, testimony.TestimonerPreviousPosition, testimony.TestimonerPhoto, testimony.TestimonyText).Scan(&testimony.ID, &testimony.CreatedAt, &testimony.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestimonyRepository) GetAllTestimonials() ([]*Testimony, error) {
	query := `
		SELECT id, testimoner_name, testimoner_current_position, testimoner_previous_position, testimoner_photo, testimony_text, created_at, updated_at, deleted_at
		FROM testimonials
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	testimonies := make([]*Testimony, 0)
	for rows.Next() {
		testimony := &Testimony{}
		if err := rows.Scan(&testimony.ID, &testimony.TestimonerName, &testimony.TestimonerCurrentPosition, &testimony.TestimonerPreviousPosition, &testimony.TestimonerPhoto, &testimony.TestimonyText, &testimony.CreatedAt, &testimony.UpdatedAt, &testimony.DeletedAt); err != nil {
			return nil, err
		}
		testimonies = append(testimonies, testimony)
	}

	return testimonies, nil
}

func (r *TestimonyRepository) GetByID(id int) (*Testimony, error) {
	query := `
		SELECT id, testimoner_name, testimoner_current_position, testimoner_previous_position, testimoner_photo, testimony_text, created_at, updated_at, deleted_at
		FROM testimonials
		WHERE id = $1 AND deleted_at IS NULL`

	testimony := &Testimony{}
	err := r.db.QueryRow(query, id).Scan(&testimony.ID, &testimony.TestimonerName, &testimony.TestimonerCurrentPosition, &testimony.TestimonerPreviousPosition, &testimony.TestimonerPhoto, &testimony.TestimonyText, &testimony.CreatedAt, &testimony.UpdatedAt, &testimony.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return testimony, nil
}

func (r *TestimonyRepository) Update(testimony *Testimony) error {
	query := `
		UPDATE testimonials
		SET testimoner_name = $1, testimoner_current_position = $2, testimoner_previous_position = $3, testimoner_photo = $4, testimony_text = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6`

	_, err := r.db.Exec(query, testimony.TestimonerName, testimony.TestimonerCurrentPosition, testimony.TestimonerPreviousPosition, testimony.TestimonerPhoto, testimony.TestimonyText, testimony.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestimonyRepository) Delete(id int) error {
	query := `
		UPDATE testimonials
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
