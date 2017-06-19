package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	conn              *sql.DB
	connStr           string
	maxIdle, maxConns int
}

func NewMysql() *Mysql {
	mysql := &Mysql{}
	return mysql
}

func (m *Mysql) Open(dbConn string, maxIdle, maxConns int) error {
	m.connStr = dbConn
	m.maxIdle = maxIdle
	m.maxConns = maxConns
	conn, err := sql.Open("mysql", dbConn)
	if err != nil {
		return err
	}
	conn.SetMaxIdleConns(maxIdle)
	conn.SetMaxOpenConns(maxConns)
	m.conn = conn
	return nil
}

func (m *Mysql) OpenOne(dbConn string) error {
	m.connStr = dbConn
	conn, err := sql.Open("mysql", dbConn)
	if err != nil {
		return err
	}
	m.conn = conn
	return nil
}

func (m *Mysql) Close() error {
	err := m.conn.Close()
	return err
}

func (m *Mysql) Insert(sql string, args ...interface{}) (int64, error) {
	stmt, err := m.conn.Prepare(sql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}
