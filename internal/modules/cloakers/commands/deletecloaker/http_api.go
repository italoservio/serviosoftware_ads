package deletecloaker

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type DeleteCloakerHttpAPI struct {
	validate              validator.Validate
	DeleteCloakerByIDCmd *DeleteCloakerByIDCmd
}

func NewDeleteCloakerHttpAPI(
	validate *validator.Validate,
	deleteCloakerByIDCmd *DeleteCloakerByIDCmd,
) *DeleteCloakerHttpAPI {
	return &DeleteCloakerHttpAPI{
		validate:              *validate,
		DeleteCloakerByIDCmd: deleteCloakerByIDCmd,
	}
}

func (d *DeleteCloakerHttpAPI) DeleteCloakerByID(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*jwt.Session)

	pathParameters := mux.Vars(r)
	cloakerID := pathParameters["cloakerId"]

	payload := DeleteCloakerByIDCmdInput{
		ID: cloakerID,
	}

	if err := d.validate.Struct(payload); err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	if err := d.DeleteCloakerByIDCmd.Exec(session, &payload); err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
