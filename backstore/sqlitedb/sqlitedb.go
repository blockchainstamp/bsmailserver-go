package sqlitedb

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-smtp"
)

type SqliteDB struct {
}

func (sd *SqliteDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (sd *SqliteDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewSqliteDB() *SqliteDB {
	sd := &SqliteDB{}
	return sd
}
func InitDefaultData(homeDir string) (err error) {
	err = util.TouchDir(homeDir)
	if err != nil {
		return err
	}
	return nil
}
