package getcloaker

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type GetCloakerHttpAPI struct {
	validate          validator.Validate
	GetCloakerByIDCmd *GetCloakerByIDCmd
}

func NewGetCloakerHttpAPI(
	validate *validator.Validate,
	getCloakerByIDCmd *GetCloakerByIDCmd,
) *GetCloakerHttpAPI {
	return &GetCloakerHttpAPI{
		validate:          *validate,
		GetCloakerByIDCmd: getCloakerByIDCmd,
	}
}

func (u *GetCloakerHttpAPI) GetCloakerByID(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*jwt.Session)

	pathParameters := mux.Vars(r)
	cloakerID := pathParameters["cloakerId"]

	body := r.Body
	defer body.Close()

	var payload GetCloakerByIDCmdInput
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

	cloaker, err := u.GetCloakerByIDCmd.Exec(session, &payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(cloaker)

	w.Write(response)
}
