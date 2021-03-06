package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lordking/blaster/common"
	"github.com/lordking/blaster/database"
)

var _ database.Database = (*MySQL)(nil)

type (
	Config struct {
		Host         string `json:"host" env:"MYSQL_HOST"`
		Port         string `json:"port" env:"MYSQL_PORT"`
		Username     string `json:"username" env:"MYSQL_USERNAME"`
		Password     string `json:"password" env:"MYSQL_PASSWORD"`
		Database     string `json:"database" env:"MYSQL_DATABASE"`
		MaxOpenConns int    `json:"maxOpenConns" env:"MYSQL_MAXOPENCONNS"`
		MaxIdleConns int    `json:"maxIdleConns" env:"MYSQL_MAXIDLECONNS"`
	}

	MySQL struct {
		Config     *Config
		Connection *sql.DB
	}
)

func (m *MySQL) NewConfig() interface{} {
	m.Config = &Config{}
	return m.Config
}

func (m *MySQL) ValidateBefore() error {

	if m.Config.Host == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `host` in config file and `MYSQL_HOST` in env")
	}

	if m.Config.Port == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `port` in config file and `MYSQL_PORT` in env")
	}

	if m.Config.Username == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `username` in config file and `MYSQL_USERNAME` in env")
	}

	if m.Config.Database == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `database` in config file and `MYSQL_DATABASE` in env")
	}

	if m.Config.MaxOpenConns < 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `maxOpenConns` in config file and `MYSQL_MAXOPENCONNS` in env")
	}

	if m.Config.MaxIdleConns < 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `maxIdleConns` in config file and `MYSQL_MAXIDLECONNS` in env")
	}

	return nil
}

func (m *MySQL) Connect() error {

	var (
		db  *sql.DB
		err error
	)

	if db, err = sql.Open("mysql", m.url()); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	db.SetMaxOpenConns(m.Config.MaxOpenConns) //最大打开的连接数，默认值为0表示不限制
	db.SetMaxIdleConns(m.Config.MaxIdleConns) //闲置的连接数

	if err = db.Ping(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	m.Connection = db

	return nil
}

func (m *MySQL) GetConnection() interface{} {
	return m.Connection
}

func (m *MySQL) Close() error {
	if err := m.Connection.Close(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	return nil
}

func (m *MySQL) url() string {
	return m.Config.Username + ":" + m.Config.Password + "@tcp(" + m.Config.Host + m.Config.Port + ")/" + m.Config.Database + "?charset=utf8"
}

func New() *MySQL {
	return &MySQL{}
}
