package jsondb

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-smtp"
)

type JsonFileDB struct {
}

func (jd *JsonFileDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (jd *JsonFileDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
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
