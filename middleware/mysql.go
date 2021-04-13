package middleware

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLClient() (*sql.DB, error) {
	return sql.Open("mysql", "zhangxu:fengzhong@/test")
}
