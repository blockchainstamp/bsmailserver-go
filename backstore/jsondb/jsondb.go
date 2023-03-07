package jsondb

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type JsonFileDB struct {
}

func (j JsonFileDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewJsonDB() *JsonFileDB {
	jf := &JsonFileDB{}
	return jf
}

func InitDefaultData(jsonHome string) (err error) {
	err = util.TouchDir(jsonHome)
	if err != nil {
		return err
	}
	return nil
}
