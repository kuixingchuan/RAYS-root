package req

import (
	"errors"
	"fmt"
	"mayfly-go/pkg/cache"
	"mayfly-go/pkg/utils/stringx"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtConf jwt配置
type JwtConf struct {
	Key                    string
	ExpireTime             uint64 // 过期时间，单位分钟
	RefreshTokenExpireTime uint64 // 刷新token的过期时间，单位分钟
}

// 默认jwt配置
var jwtConf = JwtConf{
	Key:                    stringx.RandUUID(),
	ExpireTime:             60,
	RefreshTokenExpireTime: 360,
}

// SetJwtConf 设置jwt配置
func SetJwtConf(conf JwtConf) {
	jwtConf = conf
}

// 创建用户token
func CreateToken(userId uint64, username string) (accessToken string, refreshToken string, err error) {
	tokenEpochMu.Lock()
	defer tokenEpochMu.Unlock()
	epoch := cache.GetStr(epochKey(userId))
	if epoch == "" {
		epoch = stringx.RandUUID()
	}
	if err = cache.Set(epochKey(userId), epoch, time.Minute*time.Duration(jwtConf.RefreshTokenExpireTime)); err != nil {
		return
	}
	now := time.Now()

	// 带权限创建令牌
	// 设置有效期，过期需要重新登录获取token
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": username,
		"epoch":    epoch,
		"kind":     "access",
		"exp":      now.Add(time.Minute * time.Duration(jwtConf.ExpireTime)).Unix(),
	})

	// refresh token
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": username,
		"epoch":    epoch,
		"kind":     "refresh",
		"exp":      now.Add(time.Minute * time.Duration(jwtConf.RefreshTokenExpireTime)).Unix(),
	})

	// 使用自定义字符串加密 and get the complete encoded token as a string
	accessToken, err = accessJwt.SignedString([]byte(jwtConf.Key))
	if err != nil {
		return
	}

	refreshToken, err = refreshJwt.SignedString([]byte(jwtConf.Key))
	return
}

// A per-user epoch revokes both access and refresh tokens on logout.
var tokenEpochMu sync.Mutex

func epochKey(userId uint64) string { return fmt.Sprintf("auth:token-epoch:%d", userId) }
func RevokeTokens(userId uint64) error {
	tokenEpochMu.Lock()
	defer tokenEpochMu.Unlock()
	return cache.Set(epochKey(userId), stringx.RandUUID(), time.Minute*time.Duration(jwtConf.RefreshTokenExpireTime))
}
func ParseAccessToken(value string) (uint64, string, error) { return parseTokenKind(value, "access") }
func ParseToken(value string) (uint64, string, error)       { return parseTokenKind(value, "refresh") }
func parseTokenKind(value, kind string) (uint64, string, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (any, error) { return []byte(jwtConf.Key), nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || token == nil || !token.Valid {
		return 0, "", errors.New("token invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("token invalid")
	}
	id, ok := claims["id"].(float64)
	if !ok || id <= 0 {
		return 0, "", errors.New("token invalid")
	}
	userId := uint64(id)
	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", errors.New("token invalid")
	}
	epoch, ok := claims["epoch"].(string)
	if !ok || epoch == "" || claims["kind"] != kind || cache.GetStr(epochKey(userId)) != epoch {
		return 0, "", errors.New("session revoked")
	}
	return userId, username, nil
}
