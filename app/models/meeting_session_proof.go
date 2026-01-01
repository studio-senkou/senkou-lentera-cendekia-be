package models

import (
	"database/sql"
	"time"
)

type MeetingSessionProof struct {
	ID               uint       `json:"id"`
	MeetingID        uint       `json:"meeting_id"`
	StudentProof     string     `json:"student_proof"`
	StudentSignature string     `json:"student_signature"`
	MentorProof      string     `json:"mentor_proof"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

type MeetingSessionProofRepository struct {
	db *sql.DB
}

func NewMeetingSessionProofRepository(db *sql.DB) *MeetingSessionProofRepository {
	return &MeetingSessionProofRepository{
		db: db,
	}
}
