package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CloakerConfiguration struct {
	AllowOnlyMobile bool `bson:"allowOnlyMobile" json:"allowOnlyMobile"`
}

type Cloaker struct {
	ID        bson.ObjectID        `bson:"_id,omitempty" json:"id,omitempty,string" validate:"omitempty,omitnil,mongodb"`
	UserID    bson.ObjectID        `bson:"userId" json:"userId" validate:"required,mongodb"`
	URL       string               `bson:"url" json:"url" validate:"required,url"`
	WhiteURL  string               `bson:"whiteUrl" json:"whiteUrl" validate:"required,url"`
	BlackURL  string               `bson:"blackUrl" json:"blackUrl" validate:"required,url"`
	Config    CloakerConfiguration `bson:"config" json:"config" validate:"required"`
	IsActive  bool                 `bson:"isActive" json:"isActive"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time           `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

func (c *Cloaker) StringID() string {
	return c.ID.Hex()
}

func (c *Cloaker) StringUserID() string {
	return c.UserID.Hex()
}
