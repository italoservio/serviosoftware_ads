package deps

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/clients"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/cloakerredirect"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/createcloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/deletecloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/getcloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/listcloakers"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/commands/updatecloaker"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/cache"
	"github.com/italoservio/serviosoftware_ads/pkg/db"
	"github.com/italoservio/serviosoftware_ads/pkg/env"
	"github.com/italoservio/serviosoftware_ads/pkg/validation"
)

type Container struct {
	DB              db.DB
	Env             env.Env
	Validator       validator.Validate
	Cache           cache.CacheRepository
	CloakersHttpAPI commands.CloakersHttpAPI
}

func NewContainer(envVars *env.Env) *Container {
	dbConn, err := db.NewDB((*envVars).MONGODB_URI)
	if err != nil {
		panic(err)
	}

	validator, err := validation.NewValidator()
	if err != nil {
		panic(err)
	}

	cacheRepository := cache.NewGoCacheRepository(5*time.Minute, 10*time.Minute)

	cloakersRepository := repos.NewMongoCloakerRepository(dbConn)
	ipLookupRepository := repos.NewMongoIPLookupRepository(dbConn)
	netifyClient := clients.NewNetifyClient(*envVars)

	getCloakerByIDCmd := getcloaker.NewGetCloakerByIDCmd(cloakersRepository)
	createCloakerCmd := createcloaker.NewCreateCloakerCmd(cloakersRepository, *envVars)
	deleteCloakerCmd := deletecloaker.NewDeleteCloakerByIDCmd(cloakersRepository)
	listCloakersCmd := listcloakers.NewListCloakerCmd(cloakersRepository)
	updateCloakerCmd := updatecloaker.NewUpdateCloakerByIDCmd(cloakersRepository)
	redirectCloakerCmd := cloakerredirect.NewRedirectCloakerCmd(ipLookupRepository, cloakersRepository, *netifyClient)

	cloakersHttpAPI := commands.NewCloakersHttpAPI(
		getcloaker.NewGetCloakerHttpAPI(validator, getCloakerByIDCmd),
		createcloaker.NewCreateCloakerHttpAPI(validator, createCloakerCmd),
		deletecloaker.NewDeleteCloakerHttpAPI(validator, deleteCloakerCmd),
		listcloakers.NewListCloakersHttpAPI(validator, listCloakersCmd),
		updatecloaker.NewUpdateCloakerHttpAPI(validator, updateCloakerCmd),
		cloakerredirect.NewRedirectCloakerHttpAPI(validator, redirectCloakerCmd, cacheRepository),
	)

	return &Container{
		DB:              *dbConn,
		Env:             *envVars,
		Validator:       *validator,
		Cache:           cacheRepository,
		CloakersHttpAPI: *cloakersHttpAPI,
	}
}
