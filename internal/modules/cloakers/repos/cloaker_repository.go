package repos

import "github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"

type CloakerRepository interface {
	GenerateID() string
	CreateWithID(id string, cloaker *models.Cloaker) (*models.Cloaker, error)
	Create(cloaker *models.Cloaker) (*models.Cloaker, error)
	GetByID(id string) (*models.Cloaker, error)
	List(input *ListCloakersInput) (*ListCloakersOutput, error)
	UpdateByID(id string, cloaker *models.Cloaker) (*models.Cloaker, error)
	DeleteByID(id string) error
}

type ListCloakersInput struct {
	Limit  int64   `validate:"min=1,max=100"`
	Page   int64   `validate:"min=1"`
	SortBy *string `validate:"omitempty,omitnil,oneof=createdAt"`
	Order  *string `validate:"omitempty,omitnil,oneofsortorder"`
	UserID *string `validate:"omitempty,omitnil,mongodb"`
}

type ListCloakersOutput struct {
	Total int64            `bson:"total" json:"total"`
	Items []models.Cloaker `bson:"items" json:"items"`
}
