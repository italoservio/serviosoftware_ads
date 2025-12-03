package listcloakers

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type ListCloakerCmd struct {
	repo repos.CloakerRepository
}

func NewListCloakerCmd(repo repos.CloakerRepository) *ListCloakerCmd {
	return &ListCloakerCmd{repo}
}

type ListCloakerCmdInput struct {
	*repos.ListCloakersInput
}

type ListCloakerCmdOutput struct {
	*repos.ListCloakersOutput
}

func (c *ListCloakerCmd) Exec(
	session *jwt.Session,
	input *ListCloakerCmdInput,
) (*ListCloakerCmdOutput, error) {
	if session == nil {
		return nil, exception.NewForbiddenException()
	}

	if !session.IsAdmin {
		input.UserID = &session.UserID
	}

	if input.SortBy == nil {
		createdAt := "createdAt"
		input.SortBy = &createdAt
	}

	if input.Order == nil {
		desc := "desc"
		input.Order = &desc
	}

	output, err := c.repo.List(input.ListCloakersInput)
	if err != nil {
		return nil, err
	}

	return &ListCloakerCmdOutput{output}, nil
}
