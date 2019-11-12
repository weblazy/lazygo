package auth

import (
	"fmt"
	"lazygo/core/aescode"
	"lazygo/core/consistenthash/unsafehash"
	"lazygo/core/database/redis"
	"strconv"
	"strings"
	"time"
)

const maxCount = 1000

type (
	AuthConfig struct {
		Auth          bool
		RedisNodeList []RedisNode
		MaxCount      uint32
	}
	RedisNode struct {
		RedisConf redis.RedisConf
		Position  uint32
	}
	Auth struct {
		Auth      bool
		cHashRing *unsafehash.Consistent
	}
)

var (
	AuthManager   = new(Auth)
	TokenNotFound = fmt.Errorf("token not found")
	TokenInValid  = fmt.Errorf("token is invalid")
	prefix        = "token#"
)

func InitAuth(conf AuthConfig) error {
	cHashRing := unsafehash.NewConsistent(conf.MaxCount)
	for _, value := range conf.RedisNodeList {
		if err := value.RedisConf.Validate(); err != nil {
			return err
		}
		redis := redis.NewRedis(value.RedisConf.Host, value.RedisConf.Type, value.RedisConf.Pass)
		cHashRing.Add(unsafehash.NewNode(value.RedisConf.Host, value.Position, redis))
	}
	AuthManager.cHashRing = cHashRing
	return nil
}

func (auth *Auth) Validate(token string) (string, error) {
	tokenRaw, err := aescode.AesDecrypt(token, aescode.Key)
	arr := strings.Split(tokenRaw, ":")
	if len(arr) < 2 {
		return "", TokenNotFound
	}
	if err != nil {
		return "", TokenNotFound
	}
	node := auth.cHashRing.Get(arr[0])
	value, err := node.Extra.(*redis.Redis).Get(prefix + arr[0])
	if value != token {
		return "", TokenInValid
	}
	if err != nil {
		return "", TokenNotFound
	}
	return value, nil
}

func (auth *Auth) Add(id, value string) error {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	token := id + ":" + nowStr
	node := auth.cHashRing.Get(id)
	err := node.Extra.(*redis.Redis).Set(prefix+id, token)
	if err != nil {
		return err
	}
	return nil
}
