package tokens

import (
	"time"

	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/security"
	"github.com/templwind/sass-starter/internal/svc"

	"github.com/golang-jwt/jwt/v4"
)

const AccountTokenPrefix = "account"

// NewAccountToken generates and returns a new auth record authentication token.
func NewAccountToken(svcCtx *svc.ServiceContext, userAccount *models.UserAccount) (string, error) {
	duration, _ := time.ParseDuration(svcCtx.Config.Auth.TokenDuration)
	expirationTime := time.Now().Add(duration)

	return security.NewJWT(
		jwt.MapClaims{
			"id": userAccount.AccountID,
		},
		(AccountTokenPrefix + svcCtx.Config.Auth.TokenSecret),
		expirationTime.Unix(),
	)
}

// NewUserAuthToken generates and returns a new auth record authentication token.
func NewUserAuthToken(svcCtx *svc.ServiceContext, user *models.User) (string, error) {
	duration, _ := time.ParseDuration(svcCtx.Config.Auth.TokenDuration)
	expirationTime := time.Now().Add(duration)

	return security.NewJWT(
		jwt.MapClaims{
			"id": user.ID,
		},
		(user.TokenKey + svcCtx.Config.Auth.TokenSecret),
		expirationTime.Unix(),
	)
}
