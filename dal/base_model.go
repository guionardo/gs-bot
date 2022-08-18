package dal

import "time"

type (
	// Basic model for repository without ID field
	Model struct {
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
