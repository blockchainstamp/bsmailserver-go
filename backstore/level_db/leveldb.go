package level_db

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type LevelDB struct {
}

func (l LevelDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
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
