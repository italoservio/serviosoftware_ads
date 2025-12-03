package updatecloaker

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type UpdateCloakerByIDCmd struct {
	repo repos.CloakerRepository
}

func NewUpdateCloakerByIDCmd(repo repos.CloakerRepository) *UpdateCloakerByIDCmd {
	return &UpdateCloakerByIDCmd{repo}
}

type UpdateCloakerByIDCmdInput struct {
	ID       string                       `bson:"_id,omitempty" json:"id,omitempty,string" validate:"omitempty,omitnil,mongodb"`
	URL      string                       `bson:"url" json:"url" validate:"omitempty,omitnil,url"`
	WhiteURL string                       `bson:"whiteUrl" json:"whiteUrl" validate:"omitempty,omitnil,url"`
	BlackURL string                       `bson:"blackUrl" json:"blackUrl" validate:"omitempty,omitnil,url"`
	Config   *models.CloakerConfiguration `bson:"config,omitempty" json:"config,omitempty" validate:"omitempty,omitnil"`
	IsActive *bool                        `bson:"isActive,omitempty" json:"isActive,omitempty" validate:"omitempty,omitnil"`
}

func (c *UpdateCloakerByIDCmd) Exec(
	session *jwt.Session,
	input *UpdateCloakerByIDCmdInput,
) (*models.Cloaker, error) {
	if session == nil {
		return nil, exception.NewForbiddenException()
	}

	existing, err := c.repo.GetByID(input.ID)
	if err != nil {
		return nil, exception.NewRepoException(
			"nao foi possivel buscar o cloaker, tente novamente mais tarde",
			err.Error(),
		)
	}

	if existing == nil {
		return nil, exception.NewNotFoundException("cloaker nao encontrado")
	}

	if !session.IsAdmin && existing.UserID != session.UserID {
		return nil, exception.NewForbiddenException()
	}

	updateData := &models.Cloaker{
		URL:      input.URL,
		WhiteURL: input.WhiteURL,
		BlackURL: input.BlackURL,
	}

	if input.Config != nil {
		updateData.Config = *input.Config
	}

	if input.IsActive != nil {
		updateData.IsActive = *input.IsActive
	} else {
		updateData.IsActive = existing.IsActive
	}

	updated, err := c.repo.UpdateByID(input.ID, updateData)
	if err != nil {
		return nil, exception.NewRepoException(
			"nao foi possivel atualizar o cloaker, tente novamente mais tarde",
			err.Error(),
		)
	}

	return updated, nil
}
