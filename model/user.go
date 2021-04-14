package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64
	Name     string
	Status   int
	Password string
	Regtime  string
}

func (user *User) AddUser(name string, password string) error {
	// TODO
	// 检查重名
	user.Name = name
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err == nil {
		user.Password = string(hash)
		user.Regtime = time.Now().String()[0:19]
		db, err := getMysqlClient()
		if err == nil {
			defer db.Close()
			stmt, err := db.Prepare("INSERT INTO user (name, password, regtime) VALUES(?, ?, ?)")
			if err == nil {
				defer stmt.Close()
				res, err := stmt.Exec(user.Name, user.Password, user.Regtime)
				if err == nil {
					user.Id, err = res.LastInsertId()
					if err == nil {
						return nil
					}
				}
			}
		}
	}
	return err
}
