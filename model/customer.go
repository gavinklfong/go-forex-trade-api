package model

import "time"

type Customer struct {
	ID        string
	Name      string
	Tier      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
