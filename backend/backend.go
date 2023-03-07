package backend

import (
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	_inst       StoreBackend = nil
	_once       sync.Once
	_backendLog = logrus.WithFields(logrus.Fields{
		"mode": "backend service",
	})
)

type StoreBackend interface {
}

func Inst() StoreBackend {
	conf := cfg.CurBackendConf()
	_once.Do(func() {
		switch conf.DBType {
		case cfg.DBTypMem:
			_inst = newMemDB()
		case cfg.DBTypJson:
			_inst = newJsonDB()
		case cfg.DBTypLevelDB:
			_inst = newLevelDB()
		case cfg.DBTypSqlite:
			_inst = newSqliteDB()
		}
	})
	return _inst
}
