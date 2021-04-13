package middleware

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"richangfan/forum/model"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Register string `json:"register"`
	Login    string `json:"login"`
	Logout   string `json:"logout"`
	Token    string `json:"token"`
}

type Token struct {
	Key   string
	Sum   string
	Value string
}

const TOKEN_MD5_SALT = "FOI2JF28039joijo"

const USER_CACHE_PREFIX = "user_"

func CheckLogin(c *gin.Context) User {
	var token Token
	token.Value = c.Query("token")
	if token.Value != "" {
		arr := strings.Split(token.Value, "_")
		if len(arr) == 2 {
			token.Key = arr[0]
			token.Sum = arr[1]
			ctx := context.Background()
			rdb := GetRedisClient()
			cache, err := rdb.Get(ctx, USER_CACHE_PREFIX+token.Key).Result()
			if err == nil {
				var user User
				err = json.Unmarshal([]byte(cache), &user)
				if err == nil {
					if token.Value == user.Token {
						return user
					}
				}
			}
		}
	}
	c.AbortWithStatus(http.StatusUnauthorized)
	return User{}
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
	err := validateInvitationCode(user.Code)
	if err == nil {
		var usermodel model.User
		err := usermodel.AddUser(user.Name, user.Password)
		if err == nil {
			user.Id = int(usermodel.Id)
			user.Register = usermodel.Regtime
			user.Login = usermodel.Regtime
			user.Token = generateToken(*user).Value
			value, err := json.Marshal(user)
			if err == nil {
				rdb := GetRedisClient()
				_, err = rdb.Set(context.Background(), USER_CACHE_PREFIX+strconv.Itoa(user.Id), value, 0).Result()
				if err == nil {
					return nil
				}
			}
		}
	}
	return err
}

func generateToken(user User) Token {
	var token Token
	token.Key = strconv.Itoa(user.Id)
	h := md5.New()
	h.Write([]byte(user.Name + user.Register))
	sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	h.Reset()
	h.Write([]byte(TOKEN_MD5_SALT + sum))
	token.Sum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	token.Value = token.Key + token.Sum
	return token
}

func validateInvitationCode(code string) error {
	if code != "invitecode123456" {
		return errors.New("邀请码错误")
	}
	return nil
}
