package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type IPLookupGeoLocation struct {
	Country   string `bson:"country" json:"country" validate:"required"`
	Continent string `bson:"continent" json:"continent" validate:"required"`
}

type IPLookup struct {
	ID            bson.ObjectID       `bson:"_id,omitempty" json:"id,omitempty,string" validate:"omitempty,omitnil,mongodb"`
	IP            string              `bson:"ip" json:"ip" validate:"required,ip"`
	IPPattern     string              `bson:"ipPattern" json:"ipPattern" validate:"required"`
	IsBlacklisted bool                `bson:"isBlacklisted" json:"isBlacklisted" validate:"required"`
	Applications  []string            `bson:"applications" json:"applications" validate:"required,min=1,dive,required"`
	SharedScore   float64             `bson:"sharedScore" json:"sharedScore" validate:"required,min=0"`
	GeoLocation   IPLookupGeoLocation `bson:"geoLocation" json:"geoLocation" validate:"required,dive"`
	AccessCount   int                 `bson:"accessCount" json:"accessCount" validate:"required,min=0"`
	CreatedAt     time.Time           `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt     time.Time           `bson:"updatedAt" json:"updatedAt" validate:"required"`
	DeletedAt     time.Time           `bson:"deletedAt" json:"deletedAt,omitempty"`
}

func (c *IPLookup) StringID() string {
	return c.ID.Hex()
}
