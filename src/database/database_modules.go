package database

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type LinkReference struct {
	Id       uuid.UUID `json:"id,omitempty"`
	LinkHash string    `json:"link_hash,omitempty"`
	Link     string    `json:"link,omitempty"`
}
