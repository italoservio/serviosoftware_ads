package repos

import "github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"

type IPLookupRepository interface {
	GetByIPPattern(pattern string) (*models.IPLookup, error)
	Create(ipLookup *models.IPLookup) (*models.IPLookup, error)
	IncrementAccessCount(id string) error
}
