package main

import (
	"maelstrom-echo/internal/domain"
	"maelstrom-echo/internal/service/id"
)


func UniqueIds() domain.ID {
	generator := id.NewGenerator()
	return  generator.Generate()
}