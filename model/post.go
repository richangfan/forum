package model

import (
	"errors"
	"richangfan/forum/middleware"
	"richangfan/forum/tool"
	"unicode/utf8"
)

type Post struct {
	Id      int64  `json:"id"`
	UserId  int64  `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int    `json:"status"`
	Created string `json:"created"`
}

func (post Post) AddPost() error {
	if post.UserId == 0 {
		return errors.New("未填写用户ID")
	}
	if post.Title == "" {
		return errors.New("未填写标题")
	}
	if post.Content == "" {
		return errors.New("未填写内容")
	}
	tlen := utf8.RuneCountInString(post.Title)
	if tlen > 50 {
		return errors.New("标题过长")
	}
	clen := utf8.RuneCountInString(post.Content)
	if clen > 10000 {
		return errors.New("内容过长")
	}
	db, err := middleware.GetMysqlClient()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO post (user_id, title, content, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.UserId, post.Title, post.Content, tool.GetCurrentDateTime())
	return err
}

func (post Post) GetTotal() (int64, error) {
	db, err := middleware.GetMysqlClient()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT count(*) AS total FROM post")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var total int64
	err = stmt.QueryRow().Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (post Post) GetList(start int64, end int64) ([]Post, error) {
	return nil, nil
}
