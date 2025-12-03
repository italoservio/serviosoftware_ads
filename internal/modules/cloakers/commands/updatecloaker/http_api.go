package updatecloaker

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type UpdateCloakerHttpAPI struct {
	validate              validator.Validate
	UpdateCloakerByIDCmd *UpdateCloakerByIDCmd
}

func NewUpdateCloakerHttpAPI(
	validate *validator.Validate,
	updateCloakerByIDCmd *UpdateCloakerByIDCmd,
) *UpdateCloakerHttpAPI {
	return &UpdateCloakerHttpAPI{
		validate:              *validate,
		UpdateCloakerByIDCmd: updateCloakerByIDCmd,
	}
}

func (u *UpdateCloakerHttpAPI) UpdateCloakerByID(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*jwt.Session)

	pathParameters := mux.Vars(r)
	cloakerID := pathParameters["cloakerId"]

	body := r.Body
	defer body.Close()

	var payload UpdateCloakerByIDCmdInput
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		exception.NewPayloadException().WriteJSON(w)
		return
	}

	payload.ID = cloakerID

	err = u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	cloaker, err := u.UpdateCloakerByIDCmd.Exec(session, &payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(cloaker)

	w.Write(response)
}
