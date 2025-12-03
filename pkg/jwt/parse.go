package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Parse(secret string, tkn string) (*Session, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	session := &Session{}

	if userID, ok := claims["userId"].(string); ok {
		session.UserID = userID
	}

	if roles, ok := claims["roles"].([]any); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				session.Roles = append(session.Roles, roleStr)
			}
		}
	}

	if startedAt, err := time.Parse(time.RFC3339, claims["startedAt"].(string)); err == nil {
		session.StartedAt = startedAt
	}

	if expired, ok := claims["exp"].(int64); ok {
		session.Expired = time.Unix(expired, 0).Before(time.Now())
	}

	return session, nil
}
