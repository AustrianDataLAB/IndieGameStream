package mocks

import (
	"api/models"
	"api/shared"
	"github.com/google/uuid"
)

func GameMock(identifier string) *models.Game {
	return &models.Game{
		ID:              uuid.New(),
		Title:           "Title_" + identifier,
		StorageLocation: "Storage_Location_" + identifier,
		Status:          shared.Status_New,
		Url:             "Url_" + identifier,
		Owner:           "Owner_" + identifier,
	}
}
