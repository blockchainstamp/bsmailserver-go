package level_db

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-smtp"
)

type LevelDB struct {
}

func (ld *LevelDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (ld *LevelDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewLevelDB() *LevelDB {
	ld := &LevelDB{}
	return ld
}
func InitDefaultData(homeDir string) (err error) {
	err = util.TouchDir(homeDir)
	if err != nil {
		return err
	}
	return nil
}
