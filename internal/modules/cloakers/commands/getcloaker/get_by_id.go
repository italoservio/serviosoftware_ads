package getcloaker

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type GetCloakerByIDCmd struct {
	repo repos.CloakerRepository
}

func NewGetCloakerByIDCmd(repo repos.CloakerRepository) *GetCloakerByIDCmd {
	return &GetCloakerByIDCmd{repo}
}

type GetCloakerByIDCmdInput struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty,string" validate:"omitempty,omitnil,mongodb"`
}

func (c *GetCloakerByIDCmd) Exec(
	session *jwt.Session,
	input *GetCloakerByIDCmdInput,
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

	if !session.IsAdmin && existing.StringUserID() != session.UserID {
		return nil, exception.NewForbiddenException()
	}

	return existing, nil
}
