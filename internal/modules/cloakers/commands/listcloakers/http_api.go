package listcloakers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/cast"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

type ListCloakersHttpAPI struct {
	validate        validator.Validate
	ListCloakerCmd *ListCloakerCmd
}

func NewListCloakersHttpAPI(
	validate *validator.Validate,
	listCloakerCmd *ListCloakerCmd,
) *ListCloakersHttpAPI {
	return &ListCloakersHttpAPI{
		validate:        *validate,
		ListCloakerCmd: listCloakerCmd,
	}
}

func (l *ListCloakersHttpAPI) ListCloakers(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*jwt.Session)
	query := r.URL.Query()

	payload := ListCloakerCmdInput{
		ListCloakersInput: &repos.ListCloakersInput{
			Limit:  cast.StrToInt64(query.Get("limit")),
			Page:   cast.StrToInt64(query.Get("page")),
			SortBy: cast.StrToStringPtr(query.Get("sortBy")),
			Order:  cast.StrToStringPtr(query.Get("order")),
			UserID: cast.StrToStringPtr(query.Get("userId")),
		},
	}

	if err := l.validate.Struct(payload); err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	result, err := l.ListCloakerCmd.Exec(session, &payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)

	w.Write(response)
}
