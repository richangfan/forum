package model

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"richangfan/forum/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const TOKEN_MD5_SALT = "FOI2JF28039joijo"

const USER_CACHE_PREFIX = "user_"

type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Password  string `json:"password"`
	Code      string `json:"code"`
	Regtime   string `json:"regtime"`
	Logintime string `json:"logintime"`
	Token     string `json:"token"`
}

type Token struct {
	Key   string
	Sum   string
	Value string
}

func GetUserByName(name string) (User, error) {
	if name == "" {
		return User{}, errors.New("未填写用户名")
	}
	db, err := middleware.GetMysqlClient()
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM user WHERE name = ?")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()
	var user User
	err = stmt.QueryRow(name).Scan(&user.Id, &user.Name, &user.Status, &user.Password, &user.Regtime)
	if err != nil {
		return User{}, err
	}
	if user.Id == 0 {
		return User{}, errors.New("找不到用户")
	}
	return user, nil
}

func GetUserByToken(c *gin.Context) User {
	var token Token
	token.Value = c.Query("token")
	if token.Value != "" {
		arr := strings.Split(token.Value, "_")
		if len(arr) == 2 {
			token.Key = arr[0]
			token.Sum = arr[1]
			ctx := context.Background()
			rdb := middleware.GetRedisClient()
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

func (user *User) Register() error {
	nlen := len(user.Name)
	if nlen <= 0 || nlen > 16 {
		return errors.New("用户名长度错误")
	}
	plen := len(user.Password)
	if plen <= 0 || plen > 64 {
		return errors.New("密码长度错误")
	}
	if user.Code != "invitecode123456" {
		return errors.New("邀请码错误")
	}
	_, err := GetUserByName(user.Name)
	if err == nil {
		return errors.New("用户名重复")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Regtime = time.Now().String()[0:19]
	db, err := middleware.GetMysqlClient()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO user (name, password, regtime) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.Password, user.Regtime)
	if err != nil {
		return err
	}
	user.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	user.Password = ""
	user.Code = ""
	user.Logintime = user.Regtime
	token, err := generateToken(*user)
	if err != nil {
		return err
	}
	user.Token = token.Value
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}
	rdb := middleware.GetRedisClient()
	_, err = rdb.Set(context.Background(), USER_CACHE_PREFIX+strconv.FormatInt(user.Id, 10), value, 0).Result()
	return err
}

func (user *User) Login() error {
	if user.Name == "" || user.Password == "" {
		return errors.New("用户名或密码错误")
	}
	usermodel, err := GetUserByName(user.Name)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(usermodel.Password), []byte(user.Password)) != nil {
		return errors.New("用户名或密码错误")
	}
	user.Id = usermodel.Id
	user.Name = usermodel.Name
	user.Status = usermodel.Status
	user.Password = ""
	user.Regtime = usermodel.Regtime
	user.Logintime = time.Now().String()[0:19]
	token, err := generateToken(*user)
	if err != nil {
		return err
	}
	user.Token = token.Value
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}
	rdb := middleware.GetRedisClient()
	_, err = rdb.Set(context.Background(), USER_CACHE_PREFIX+strconv.FormatInt(user.Id, 10), value, 0).Result()
	return err
}

func generateToken(user User) (Token, error) {
	if user.Id == 0 || user.Name == "" || user.Logintime == "" {
		return Token{}, errors.New("参数错误")
	}
	var token Token
	token.Key = strconv.FormatInt(user.Id, 10)
	h := md5.New()
	h.Write([]byte(user.Name + user.Logintime))
	sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	h.Reset()
	h.Write([]byte(TOKEN_MD5_SALT + sum))
	token.Sum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	token.Value = token.Key + "_" + token.Sum
	return token, nil
}
