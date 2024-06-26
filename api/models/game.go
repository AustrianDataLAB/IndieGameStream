package models

import (
	"api/shared"
	"github.com/google/uuid"
)

type Game struct {
	ID              uuid.UUID         `json:"id"`
	Title           string            `json:"title"`
	StorageLocation string            `json:"storageLocation"`
	Status          shared.GameStatus `json:"status"`
	Url             string            `json:"url"`
	Owner           string            `json:"owner"`
	FileName        string            `json:"fileName"`
}
