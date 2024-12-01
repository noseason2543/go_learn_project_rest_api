package auth

import (
	"errors"
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/users"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type auth struct {
	mapClaims *mapClaims
	cfg       config.IJwtConfig
}

type admin struct {
	*auth
}

type mapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IAuth interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nonShop-api",                  // create by
				Subject:   "refresh-token",                // purpose of this token
				Audience:  []string{"customer", "admin"},  // who can use
				ExpiresAt: jwtTimeRepeatAdapter(exp),      // expired at
				NotBefore: jwt.NewNumericDate(time.Now()), // token is not available until time that set
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}

}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nonShop-api",                             // create by
				Subject:   "access-token",                            // purpose of this token
				Audience:  []string{"customer", "admin"},             // who can use
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()), // expired at
				NotBefore: jwt.NewNumericDate(time.Now()),            // token is not available until time that set
				IssuedAt:  jwt.NewNumericDate(time.Now()),            // when token create
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nonShop-api",                              // create by
				Subject:   "refresh-token",                            // purpose of this token
				Audience:  []string{"customer", "admin"},              // who can use
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()), // expired at
				NotBefore: jwt.NewNumericDate(time.Now()),             // token is not available until time that set
				IssuedAt:  jwt.NewNumericDate(time.Now()),             // when token create
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuth {
	return &admin{
		auth: &auth{
			cfg: cfg,
			mapClaims: &mapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "nonShop-api",                  // create by
					Subject:   "admin-token",                  // purpose of this token
					Audience:  []string{"admin"},              // who can use
					ExpiresAt: jwtTimeDurationCal(300),        // expired at
					NotBefore: jwt.NewNumericDate(time.Now()), // token is not available until time that set
					IssuedAt:  jwt.NewNumericDate(time.Now()), // when token create
				},
			},
		},
	}
}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *admin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	claims := &mapClaims{}

	// example if need token output
	/*
		    token , err := jwt.ParseWithClaims(tokenString, Claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("signing method is invalid")
				}
				return cfg.SecretKey(), nil
			})
	*/

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("token had expired message error: %v", err)
	} else if err != nil { // if claims struct is not match error will occur
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	// if !token.Valid {
	// 	return nil, fmt.Errorf("token is invalid")
	// }

	return claims, nil
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	claims := &mapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("token had expired message error: %v", err)
	} else if err != nil { // if claims struct is not match error will occur
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	return claims, nil
}
