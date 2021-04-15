package model

import (
	"context"
	"errors"
	"fmt"
	"richangfan/forum/middleware"
	"richangfan/forum/tool"
	"unicode/utf8"
)

type Post struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
	Created  string `json:"created"`
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
	if tlen > 100 {
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

func (post Post) GetList(limit int, offset int) ([]Post, error) {
	db, err := middleware.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT post.id, post.user_id, user.name, post.title, post.content, post.status, post.created FROM post LEFT JOIN user ON post.user_id = user.id ORDER BY post.id DESC LIMIT ? OFFSET ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := make([]Post, 0)
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.Id, &p.UserId, &p.UserName, &p.Title, &p.Content, &p.Status, &p.Created)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(ps)
	return ps, nil
}
