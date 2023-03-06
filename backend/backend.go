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
		if conf.UseMemDB {
			_inst = newMemDB()
		} else {
			_inst = newSqliteDB()
		}
	})
	return _inst
}
