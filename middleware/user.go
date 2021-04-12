package middleware

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Created  string `json:"created"`
}

type Token struct {
	Key   string
	Sum   string
	Value string
}

const TOKEN_MD5_SALT = "FOI2JF28039joijo"

const TOKEN_CACHE_PREFIX = "user_token_"

const INVITE_CODE_CACHE = "invite_code"

func CheckLogin(c *gin.Context) User {
	var token Token
	token.Value = c.Query("token")
	if token.Value != "" {
		arr := strings.Split(token.Value, "_")
		if len(arr) == 2 {
			token.Key = arr[0]
			token.Sum = arr[1]
			ctx := context.Background()
			client := GetRedisClient()
			cache, err := client.Get(ctx, TOKEN_CACHE_PREFIX+token.Key).Result()
			if err == nil {
				var user User
				err = json.Unmarshal([]byte(cache), &user)
				if err == nil {
					if token.Value == (generateToken(user)).Value {
						return user
					}
				}
			}
		}
	}
	c.AbortWithStatus(http.StatusUnauthorized)
	return User{}
}

func Login(user User) (string, error) {
	text, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	token := generateToken(user)
	ctx := context.Background()
	client := GetRedisClient()
	_, err = client.Set(ctx, TOKEN_CACHE_PREFIX+token.Key, text, 0).Result()
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

func Logout(user User) error {
	ctx := context.Background()
	client := GetRedisClient()
	_, err := client.Del(ctx, TOKEN_CACHE_PREFIX+strconv.Itoa(user.Id)).Result()
	return err
}

func Register(user *User) error {
	nlen := len(user.Name)
	if nlen <= 0 || nlen > 16 {
		return errors.New("用户名长度错误")
	}
	plen := len(user.Password)
	if plen <= 0 || plen > 64 {
		return errors.New("密码长度错误")
	}
	ctx := context.Background()
	client := GetRedisClient()
	cache, err := client.Get(ctx, INVITE_CODE_CACHE).Result()
	if err != nil {
		return err
	}
	if cache != user.Code {
		return errors.New("邀请码错误")
	}
	for {
		user.Id = rand.Int()
		cache, err = client.Get(ctx, TOKEN_CACHE_PREFIX+strconv.Itoa(user.Id)).Result()
		if err != nil {
			return nil
		}
		if len(cache) == 0 {
			break
		}
	}
	// TODO
	return nil
}

func generateToken(user User) Token {
	var token Token
	token.Key = strconv.Itoa(user.Id)
	h := md5.New()
	h.Write([]byte(user.Name + user.Created))
	sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	h.Reset()
	h.Write([]byte(TOKEN_MD5_SALT + sum))
	token.Sum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	token.Value = token.Key + token.Sum
	return token
}
