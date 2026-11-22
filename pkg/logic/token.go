package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/errorx"
	"strconv"
	"time"
)

type loginClaims struct {
	jwt.RegisteredClaims
	Id           string `json:"id"`
	Version      int64  `json:"ver"`
	LegacyExpiry int64  `json:"expiresAt,omitempty"`
	LegacyIssuer string `json:"issuer,omitempty"`
}
type Auth struct {
	UID, Version int64
	Expires      time.Time
}

func (l *UserLoginLogic) issueToken(uid, version int64) (string, error) {
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}
	claims := loginClaims{Id: strconv.FormatInt(uid, 36), Version: version, RegisteredClaims: jwt.RegisteredClaims{ID: hex.EncodeToString(nonce[:]), Issuer: l.appCtx.Config().Name, ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(l.appCtx.Config().Cipher))
}
func decodeToken(raw, name, key string) (Auth, error) {
	var c loginClaims
	_, err := jwt.ParseWithClaims(raw, &c, func(*jwt.Token) (interface{}, error) { return []byte(key), nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return Auth{}, errorx.UserTokenError
	}
	expiry := time.Unix(c.LegacyExpiry, 0)
	issuer := c.LegacyIssuer
	if c.ExpiresAt != nil {
		expiry = c.ExpiresAt.Time
		issuer = c.Issuer
	}
	uid, err := strconv.ParseInt(c.Id, 36, 64)
	if err != nil || uid <= 0 || issuer != name || !expiry.After(time.Now()) {
		return Auth{}, errorx.UserTokenError
	}
	return Auth{uid, c.Version, expiry}, nil
}
func (l *UserLoginLogic) revokedKey(token string) string {
	return fmt.Sprintf("%s:revoked:%x", l.appCtx.Config().Name, sha256.Sum256([]byte(token)))
}
func (l *UserLoginLogic) ValidateToken(token string) (Auth, error) {
	auth, err := decodeToken(token, l.appCtx.Config().Name, l.appCtx.Config().Cipher)
	if err != nil {
		return Auth{}, err
	}
	n, err := l.appCtx.RedisCache().Exists(context.Background(), l.revokedKey(token)).Result()
	if err != nil {
		return Auth{}, err
	}
	v, err := l.appCtx.UserModel().TokenVersion(auth.UID)
	if err != nil {
		return Auth{}, err
	}
	if n != 0 || v != auth.Version {
		return Auth{}, errorx.UserTokenError
	}
	return auth, nil
}
func (l *UserLoginLogic) Logout(token string) error {
	auth, err := decodeToken(token, l.appCtx.Config().Name, l.appCtx.Config().Cipher)
	if err != nil {
		return err
	}
	return l.appCtx.RedisCache().Set(context.Background(), l.revokedKey(token), "1", time.Until(auth.Expires)+time.Second).Err()
}
func (l *UserLoginLogic) ChangePassword(auth Auth, old, new string) (*dto.LoginRes, error) {
	version, err := l.appCtx.UserModel().ChangePassword(auth.UID, old, new, auth.Version)
	if err != nil {
		return nil, err
	}
	token, err := l.issueToken(auth.UID, version)
	if err != nil {
		return nil, err
	}
	user, err := getUserInfo(auth.UID, l.appCtx)
	return &dto.LoginRes{User: user, Token: token}, err
}
func (l *UserLoginLogic) UpdateProfile(uid int64, req dto.UpdateProfileReq) (*dto.User, error) {
	m := l.appCtx.UserModel()
	if err := m.UpdateUser(uid, nil, req.Nickname, req.Avatar, nil, req.Sex, req.Birthday); err != nil {
		return nil, err
	}
	user, err := m.FindOne(uid)
	if err != nil {
		return nil, err
	}
	result := userModel2UserDto(user)
	return result, setUserInfoCache(result, l.appCtx)
}
