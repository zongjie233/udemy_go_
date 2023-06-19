package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB 数据库连接
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL 连接到 SQL 数据库并返回一个指向该数据库的指针。
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetConnMaxLifetime(maxDbLifetime)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetMaxOpenConns(maxOpenDbConn)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB 测试 SQL 数据库的连接是否成功
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil

}

// NewDatabase 创建新数据库
func NewDatabase(dsn string) (*sql.DB, error) {
	// 初始化一个 sql.DB 实例来代表 PostgreSQL 连接
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// 执行一个 Ping 操作，检查连接是否正常
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// 返回连接到 PostgreSQL 数据库的 sql.DB 实例
	return db, nil
}
