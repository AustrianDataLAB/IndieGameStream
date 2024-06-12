package dtos

import (
	"api/shared"
	"github.com/google/uuid"
)

type GetAllGamesResponseBody struct {
	ID     uuid.UUID         `json:"id"`
	Title  string            `json:"title"`
	Status shared.GameStatus `json:"status"`
	Url    string            `json:"url"`
}

type GetGameByIdResponseBody struct {
	ID     uuid.UUID         `json:"id"`
	Title  string            `json:"title"`
	Status shared.GameStatus `json:"status"`
	Url    string            `json:"url"`
}

type UploadGameRequestBody struct {
	Title string `json:"title"`
}
