package db

import (
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"strings"
)

const (
	EngineSqlite   = "sqlite"
	EnginePostgres = "postgresql"
	EngineMysql    = "mysql"
)

func IsSupportedEngines(engine string) bool {
	switch strings.ToLower(engine) {
	case EngineSqlite, EnginePostgres, EngineMysql:
		return true
	}
	return false
}

var instance *Database

type Database struct {
	*gorm.DB
	engine string
}

func GetInstance() *Database {
	return instance
}

func NewDB(engine string) (*Database, error) {
	if instance != nil {
		return instance, nil
	}

	var db *gorm.DB
	var err error
	switch strings.ToLower(engine) {
	case EngineSqlite:
		db, err = gorm.Open(sqlite.Open("sqlite3.db"))
	case EnginePostgres:
		// TODO
	case EngineMysql:
		// TODO
	}
	if err != nil {
		return nil, err
	}

	return &Database{
		DB:     db,
		engine: engine,
	}, nil
}

func (db *Database) MigrateTables() error {
	return db.AutoMigrate(
		&model.Auth{},
		&model.Autosave{},
		&model.Board{},
		&model.BoardFile{},
		&model.BoardGood{},
		&model.BoardNew{},
		&model.Config{},
		&model.Content{},
		&model.Faq{},
		&model.FaqMaster{},
		&model.Group{},
		&model.GroupMember{},
		&model.Login{},
		&model.Mail{},
		&model.Member{},
		&model.MemberSocialProfile{},
		&model.Memo{},
		&model.Menu{},
		&model.NewWin{},
		&model.Point{},
		&model.Poll{},
		&model.PollEtc{},
		&model.Popular{},
		&model.QaConfig{},
		&model.QaContent{},
		&model.Scrap{},
		&model.Uniqid{},
		&model.Visit{},
		&model.VisitSum{},
	)
}

func (db *Database) ListAllTables() ([]string, error) {
	var query string
	switch db.engine {
	case EngineSqlite:
		query = sqliteTablesQuery()
	case EngineMysql:
		// TODO
	case EnginePostgres:
		// TODO
	}

	var tableNames []string
	err := db.Raw(query).Scan(&tableNames).Error
	if err != nil {
		return nil, err
	}
	return tableNames, nil
}

func (db *Database) HasTable(tableName string) (bool, error) {
	var count int64
	var err error
	switch db.engine {
	case EngineSqlite:
		err = db.Table("sqlite_master").
			Where("type = 'table' AND name = ?", tableName).
			Select("count(*)").
			Scan(&count).
			Error
	case EngineMysql:
		// TODO
	case EnginePostgres:
		// TODO
	}
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
