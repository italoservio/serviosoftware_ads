package commands

import (
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/cloakerredirect"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/createcloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/deletecloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/listcloakers"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/updatecloaker"
)

type CloakersHttpAPI struct {
	CreateCloakerHttpAPI   createcloaker.CreateCloakerHttpAPI
	DeleteCloakerHttpAPI   deletecloaker.DeleteCloakerHttpAPI
	ListCloakersHttpAPI    listcloakers.ListCloakersHttpAPI
	UpdateCloakerHttpAPI   updatecloaker.UpdateCloakerHttpAPI
	RedirectCloakerHttpAPI cloakerredirect.RedirectCloakerHttpAPI
}

func NewCloakersHttpAPI(
	createCloakerHttpAPI *createcloaker.CreateCloakerHttpAPI,
	deleteCloakerHttpAPI *deletecloaker.DeleteCloakerHttpAPI,
	listCloakersHttpAPI *listcloakers.ListCloakersHttpAPI,
	updateCloakerHttpAPI *updatecloaker.UpdateCloakerHttpAPI,
	redirectCloakerHttpAPI *cloakerredirect.RedirectCloakerHttpAPI,
) *CloakersHttpAPI {
	return &CloakersHttpAPI{
		CreateCloakerHttpAPI:   *createCloakerHttpAPI,
		DeleteCloakerHttpAPI:   *deleteCloakerHttpAPI,
		ListCloakersHttpAPI:    *listCloakersHttpAPI,
		UpdateCloakerHttpAPI:   *updateCloakerHttpAPI,
		RedirectCloakerHttpAPI: *redirectCloakerHttpAPI,
	}
}
