package entity

import (
	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func FromString(s string) (ID, error) {
	uuid, err := uuid.Parse(s)
	return ID(uuid), err
}
