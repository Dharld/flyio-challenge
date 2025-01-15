package domain

type ID string

func NewID(s string) ID {
	return ID(s)
}