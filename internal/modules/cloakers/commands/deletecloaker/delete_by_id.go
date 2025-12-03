package deletecloaker

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type DeleteCloakerByIDCmd struct {
	repo repos.CloakerRepository
}

type DeleteCloakerByIDCmdInput struct {
	ID string `validate:"required,mongodb"`
}

func NewDeleteCloakerByIDCmd(repo repos.CloakerRepository) *DeleteCloakerByIDCmd {
	return &DeleteCloakerByIDCmd{repo}
}

func (c *DeleteCloakerByIDCmd) Exec(
	session *jwt.Session,
	input *DeleteCloakerByIDCmdInput,
) error {
	if session == nil {
		return exception.NewForbiddenException()
	}

	cloaker, err := c.repo.GetByID(input.ID)
	if err != nil {
		return exception.NewRepoException(
			"nao foi possivel encontrar o cloaker, tente novamente mais tarde",
			err.Error(),
		)
	}

	if cloaker == nil {
		return exception.NewNotFoundException("cloaker nao encontrado")
	}

	if !session.IsAdmin && cloaker.StringUserID() != session.UserID {
		return exception.NewForbiddenException()
	}

	err = c.repo.DeleteByID(input.ID)
	if err != nil {
		return exception.NewRepoException(
			"nao foi possivel remover o cloaker, tente novamente mais tarde",
			err.Error(),
		)
	}

	return nil
}
