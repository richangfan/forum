package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func getMysqlClient() (*sql.DB, error) {
	return sql.Open("mysql", "zhangxu:fengzhong@/forum")
}
