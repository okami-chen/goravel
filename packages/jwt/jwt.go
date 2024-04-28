package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/database"
	"github.com/spf13/cast"
	"goravel/packages/jwt/contracts"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const ctxKey = "Jwt"

type Claims struct {
	Key string `json:"key"`
	jwt.RegisteredClaims
	Expand map[string]interface{} `json:"ext"`
}

type Guard struct {
	Claims *Claims
	Token  string
}

type Guards map[string]*Guard

type Jwt struct {
	cache  cache.Cache
	config config.Config
	guard  string
	orm    orm.Orm
}

func NewJwt(guard string, cache cache.Cache, config config.Config, orm orm.Orm) *Jwt {
	return &Jwt{
		cache:  cache,
		config: config,
		guard:  guard,
		orm:    orm,
	}
}

func (a *Jwt) Guard(name string) contracts.Jwt {
	return NewJwt(name, a.cache, a.config, a.orm)
}

func (a *Jwt) User(ctx http.Context, user any, expand map[string]interface{}) error {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[a.guard] == nil {
		return ErrorParseTokenFirst
	}
	if auth[a.guard].Claims == nil {
		return ErrorParseTokenFirst
	}
	if auth[a.guard].Claims.Key == "" {
		return ErrorInvalidKey
	}
	if auth[a.guard].Token == "" {
		return ErrorTokenExpired
	}
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", "jwt", clause.PrimaryColumn.Table+":"+clause.PrimaryColumn.Name, auth[a.guard].Claims.Key)
	val := a.cache.Get(cacheKey)
	if val != nil && val != "nil" {
		return nil
	}

	if val == "nil" {
		return errors.New("record not found")
	}

	if err := a.orm.Query().FindOrFail(user, clause.Eq{Column: clause.PrimaryColumn, Value: auth[a.guard].Claims.Key}); err != nil {
		a.cache.Put(cacheKey, "nil", time.Minute*60)
		return err
	}

	a.cache.Put(cacheKey, "y", time.Minute*60)

	return nil
}

func (a *Jwt) Parse(ctx http.Context, token string) (*contracts.Payload, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	if a.cache == nil {
		return nil, errors.New("cache support is required")
	}
	if a.tokenIsDisabled(token) {
		return nil, ErrorTokenDisabled
	}

	jwtSecret := a.config.GetString("jwt.secret")
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithTimeFunc(func() time.Time {
		return carbon.Now().ToStdTime()
	}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) && tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if !ok {
				return nil, ErrorInvalidClaims
			}

			a.makeAuthContext(ctx, claims, "")

			return &contracts.Payload{
				Guard:    claims.Subject,
				Key:      claims.Key,
				ExpireAt: claims.ExpiresAt.Local(),
				IssuedAt: claims.IssuedAt.Local(),
				Expand:   claims.Expand,
			}, ErrorTokenExpired
		}

		return nil, ErrorInvalidToken
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return nil, ErrorInvalidToken
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, ErrorInvalidClaims
	}

	a.makeAuthContext(ctx, claims, token)

	return &contracts.Payload{
		Guard:    claims.Subject,
		Key:      claims.Key,
		ExpireAt: claims.ExpiresAt.Time,
		IssuedAt: claims.IssuedAt.Time,
		Expand:   claims.Expand,
	}, nil
}

func (a *Jwt) Login(ctx http.Context, user any, expand map[string]interface{}) (token string, err error) {
	id := database.GetID(user)
	if id == nil {
		return "", ErrorNoPrimaryKeyField
	}

	return a.LoginUsingID(ctx, id, expand)
}

func (a *Jwt) LoginUsingID(ctx http.Context, id any, expand map[string]interface{}) (token string, err error) {
	jwtSecret := a.config.GetString("jwt.secret")
	if jwtSecret == "" {
		return "", ErrorEmptySecret
	}

	nowTime := carbon.Now()
	ttl := a.config.GetInt("jwt.ttl")
	if ttl == 0 {
		// 100 years
		ttl = 60 * 24 * 365 * 100
	}
	expireTime := nowTime.AddMinutes(ttl).ToStdTime()
	key := cast.ToString(id)
	if key == "" {
		return "", ErrorInvalidKey
	}
	claims := Claims{
		Key: key,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime.ToStdTime()),
			Subject:   a.guard,
		},
		Expand: expand,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	a.makeAuthContext(ctx, &claims, token)

	return
}

// Refresh need parse token first.
func (a *Jwt) Refresh(ctx http.Context) (token string, err error) {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[a.guard] == nil {
		return "", ErrorParseTokenFirst
	}
	if auth[a.guard].Claims == nil {
		return "", ErrorParseTokenFirst
	}

	nowTime := carbon.Now()
	refreshTtl := a.config.GetInt("jwt.refresh_ttl")
	if refreshTtl == 0 {
		// 100 years
		refreshTtl = 60 * 24 * 365 * 100
	}

	expireTime := carbon.FromStdTime(auth[a.guard].Claims.ExpiresAt.Time).AddMinutes(refreshTtl)
	if nowTime.Gt(expireTime) {
		return "", ErrorRefreshTimeExceeded
	}

	return a.LoginUsingID(ctx, auth[a.guard].Claims.Key, auth[a.guard].Claims.Expand)
}

func (a *Jwt) Logout(ctx http.Context) error {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[a.guard] == nil || auth[a.guard].Token == "" {
		return nil
	}

	if a.cache == nil {
		return errors.New("cache support is required")
	}

	ttl := a.config.GetInt("jwt.ttl")
	if ttl == 0 {
		if ok := a.cache.Forever(getDisabledCacheKey(auth[a.guard].Token), true); !ok {
			return errors.New("cache forever failed")
		}
	} else {
		if err := a.cache.Put(getDisabledCacheKey(auth[a.guard].Token),
			true,
			time.Duration(ttl)*time.Minute,
		); err != nil {
			return err
		}
	}

	delete(auth, a.guard)
	ctx.WithValue(ctxKey, auth)

	return nil
}

func (a *Jwt) makeAuthContext(ctx http.Context, claims *Claims, token string) {
	guards, ok := ctx.Value(ctxKey).(Guards)
	if !ok {
		guards = make(Guards)
	}
	guards[a.guard] = &Guard{claims, token}
	ctx.WithValue(ctxKey, guards)
}

func (a *Jwt) tokenIsDisabled(token string) bool {
	return a.cache.GetBool(getDisabledCacheKey(token), false)
}

func getDisabledCacheKey(token string) string {
	return "jwt:disabled:" + token
}
