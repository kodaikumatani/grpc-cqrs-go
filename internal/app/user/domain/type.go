package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID        ulid.ULID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
