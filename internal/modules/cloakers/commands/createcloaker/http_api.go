package createcloaker

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
)

type CreateCloakerHttpAPI struct {
	validate         validator.Validate
	CreateCloakerCmd *CreateCloakerCmd
}

func NewCreateCloakerHttpAPI(
	validate *validator.Validate,
	createCloakerCmd *CreateCloakerCmd,
) *CreateCloakerHttpAPI {
	return &CreateCloakerHttpAPI{
		validate:         *validate,
		CreateCloakerCmd: createCloakerCmd,
	}
}

func (c *CreateCloakerHttpAPI) CreateCloaker(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	var payload CreateCloakerInput
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		exception.NewPayloadException().WriteJSON(w)
		return
	}

	if err := c.validate.Struct(payload); err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	cloaker, err := c.CreateCloakerCmd.Exec(&payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, _ := json.Marshal(cloaker)
	w.Write(response)
}
