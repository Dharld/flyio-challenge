package id

import (
	"maelstrom-echo/internal/domain"

	"github.com/google/uuid"
)

type Generator struct {}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate() domain.ID {
	id := uuid.New()
	return domain.NewID(id.String())
}