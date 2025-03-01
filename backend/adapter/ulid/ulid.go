package ulid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

type IULID interface {
	GenerateID() string
}

type ULID struct{}

func NewULID() IULID {
	return &ULID{}
}

func (u *ULID) GenerateID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
