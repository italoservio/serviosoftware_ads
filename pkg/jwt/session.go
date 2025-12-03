package jwt

import "time"

type Session struct {
	UserID    string    `bson:"userId" json:"userId" validate:"required,mongodb"`
	Roles     []string  `bson:"roles" json:"roles" validate:"required,min=1,dive,required,oneofrole"`
	StartedAt time.Time `bson:"startedAt" json:"startedAt" validate:"required"`
	IsAdmin   bool      `bson:"isAdmin" json:"isAdmin" validate:"required"`
	Expired   bool      `bson:"expired" json:"expired" validate:"required"`
}
