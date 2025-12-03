package commands

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/cloakerredirect"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/createcloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/deletecloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/getcloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/listcloakers"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/updatecloaker"
)

type CloakersHttpAPI struct {
	GetCloakerByIDHttpAPI  getcloaker.GetCloakerHttpAPI
	CreateCloakerHttpAPI   createcloaker.CreateCloakerHttpAPI
	DeleteCloakerHttpAPI   deletecloaker.DeleteCloakerHttpAPI
	ListCloakersHttpAPI    listcloakers.ListCloakersHttpAPI
	UpdateCloakerHttpAPI   updatecloaker.UpdateCloakerHttpAPI
	RedirectCloakerHttpAPI cloakerredirect.RedirectCloakerHttpAPI
}

func NewCloakersHttpAPI(
	getCloakerByIDHttpAPI *getcloaker.GetCloakerHttpAPI,
	createCloakerHttpAPI *createcloaker.CreateCloakerHttpAPI,
	deleteCloakerHttpAPI *deletecloaker.DeleteCloakerHttpAPI,
	listCloakersHttpAPI *listcloakers.ListCloakersHttpAPI,
	updateCloakerHttpAPI *updatecloaker.UpdateCloakerHttpAPI,
	redirectCloakerHttpAPI *cloakerredirect.RedirectCloakerHttpAPI,
) *CloakersHttpAPI {
	return &CloakersHttpAPI{
		GetCloakerByIDHttpAPI:  *getCloakerByIDHttpAPI,
		CreateCloakerHttpAPI:   *createCloakerHttpAPI,
		DeleteCloakerHttpAPI:   *deleteCloakerHttpAPI,
		ListCloakersHttpAPI:    *listCloakersHttpAPI,
		UpdateCloakerHttpAPI:   *updateCloakerHttpAPI,
		RedirectCloakerHttpAPI: *redirectCloakerHttpAPI,
	}
}
