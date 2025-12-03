package createcloaker

import (
	"encoding/base64"

	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/env"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
)

type CreateCloakerCmd struct {
	repo repos.CloakerRepository
	env  env.Env
}

func NewCreateCloakerCmd(repo repos.CloakerRepository, e env.Env) *CreateCloakerCmd {
	return &CreateCloakerCmd{repo: repo, env: e}
}

type CreateCloakerInput struct {
	UserID   string                      `bson:"userId" json:"userId" validate:"required,mongodb"`
	WhiteURL string                      `bson:"whiteUrl" json:"whiteUrl" validate:"required,url"`
	BlackURL string                      `bson:"blackUrl" json:"blackUrl" validate:"required,url"`
	Config   models.CloakerConfiguration `bson:"config" json:"config" validate:"required"`
}

func (c *CreateCloakerCmd) Exec(input *CreateCloakerInput) (*models.Cloaker, error) {
	generatedID := c.repo.GenerateID()
	encodedID := base64.URLEncoding.EncodeToString([]byte(generatedID))
	url := c.env.SERVER_BASE_URL + "/r/" + encodedID

	cloaker := &models.Cloaker{
		UserID:   c.repo.ToID(input.UserID),
		URL:      url,
		WhiteURL: input.WhiteURL,
		BlackURL: input.BlackURL,
		Config:   input.Config,
	}

	cloaker, err := c.repo.CreateWithID(generatedID, cloaker)
	if err != nil {
		err := exception.NewRepoException(
			"nao foi possivel criar o cloaker, tente novamente mais tarde",
			err.Error(),
		)
		return nil, err
	}

	return cloaker, nil
}
