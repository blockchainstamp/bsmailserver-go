package backstore

import (
	"github.com/blockchainstamp/bsmailserver-go/backstore/jsondb"
	"github.com/blockchainstamp/bsmailserver-go/backstore/level_db"
	"github.com/blockchainstamp/bsmailserver-go/backstore/memory"
	"github.com/blockchainstamp/bsmailserver-go/backstore/sqlitedb"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"sync"
)

var (
	_inst     StoreBackend = nil
	_once     sync.Once
	_storeLog = logrus.WithFields(logrus.Fields{
		"mode": "backend storage",
	})
)

type StoreBackend interface {
	backend.Backend
	smtp.Backend
}

func Inst() StoreBackend {
	conf := cfg.CurBackendConf()
	_once.Do(func() {
		switch conf.DBType {
		case cfg.DBTypMem:
			_inst = memory.NewMemDB()
		case cfg.DBTypJson:
			_inst = jsondb.NewJsonDB()
		case cfg.DBTypLevelDB:
			_inst = level_db.NewLevelDB()
		case cfg.DBTypSqlite:
			_inst = sqlitedb.NewSqliteDB()
		}
	})
	return _inst
}

func InitDefaultDB(dbHome string) (err error) {
	if err := util.TouchDir(dbHome); err != nil {
		return err
	}
	memDBDir := path.Join(dbHome, string(filepath.Separator), util.DBMemHome)
	if err = memory.InitDefaultData(memDBDir); err != nil {
		return err
	}
	jsonDBDir := path.Join(dbHome, string(filepath.Separator), util.DBJsonHome)
	if err = jsondb.InitDefaultData(jsonDBDir); err != nil {
		return err
	}

	levelDBDir := path.Join(dbHome, string(filepath.Separator), util.DBLevelHome)
	if err = level_db.InitDefaultData(levelDBDir); err != nil {
		return err
	}

	sqliteDBDir := path.Join(dbHome, string(filepath.Separator), util.DBSqliteHome)
	if err = sqlitedb.InitDefaultData(sqliteDBDir); err != nil {
		return err
	}
	return nil
}
