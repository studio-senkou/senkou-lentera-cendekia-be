package models

import "database/sql"

type StaticAsset struct {
	ID        int    `json:"id"`
	AssetName string `json:"asset_name"`
	AssetType string `json:"asset_type"` // e.g., 'image', 'video', 'document'
	AssetURL  string `json:"asset_url"`
	AssetDescription string `json:"asset_description,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type StaticAssetRepository struct {
	db *sql.DB
}

func NewStaticAssetRepository(db *sql.DB) *StaticAssetRepository {
	return &StaticAssetRepository{db: db}
}

func (r *StaticAssetRepository) Create(asset *StaticAsset) error {
	query := `
		INSERT INTO static_assets (asset_name, asset_type, asset_url, asset_description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, asset.AssetName, asset.AssetType, asset.AssetURL, asset.AssetDescription).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)
}

func (r *StaticAssetRepository) GetAll() ([]*StaticAsset, error) {
	query := `SELECT id, asset_name, asset_type, asset_url, asset_description, created_at, updated_at FROM static_assets ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := make([]*StaticAsset, 0)
	for rows.Next() {
		var asset StaticAsset
		if err := rows.Scan(&asset.ID, &asset.AssetName, &asset.AssetType, &asset.AssetURL, &asset.AssetDescription, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (r *StaticAssetRepository) GetByID(id int) (*StaticAsset, error) {
	query := `SELECT id, asset_name, asset_type, asset_url, asset_description, created_at, updated_at FROM static_assets WHERE id = $1`
	asset := new(StaticAsset)
	err := r.db.QueryRow(query, id).Scan(&asset.ID, &asset.AssetName, &asset.AssetType, &asset.AssetURL, &asset.AssetDescription, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return asset, nil
}

func (r *StaticAssetRepository) Delete(id int) error {
	query := `DELETE FROM static_assets WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}