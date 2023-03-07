package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
)

const (
	MemDbFileName = "memDB.json"
)

var (
	_memBackLog = logrus.WithFields(logrus.Fields{
		"mode": "backend storage",
	})
	defaultUser = []string{
		"ribencong@simplenets.org",
		"fireman@simplenets.org",
	}
	defaultMailBox = map[string]*ImapMailBox{
		"INBOX": &ImapMailBox{
			Subscribed: false,
			BName:      "INBOX",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{},
				Delimiter:  "/",
				Name:       "INBOX",
			},
		},
		"草稿箱": &ImapMailBox{
			Subscribed: false,
			BName:      "草稿箱",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Drafts"},
				Delimiter:  "/",
				Name:       "草稿箱",
			},
		},
		"已发送": &ImapMailBox{
			Subscribed: false,
			BName:      "已发送",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Sent"},
				Delimiter:  "/",
				Name:       "已发送",
			},
		},
		"已删除": &ImapMailBox{
			Subscribed: false,
			BName:      "已删除",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Trash"},
				Delimiter:  "/",
				Name:       "已删除",
			},
		},
		"区块链邮票": &ImapMailBox{
			Subscribed: false,
			BName:      "区块链邮票",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{},
				Delimiter:  "/",
				Name:       "区块链邮票",
			},
		},
	}
)

type MemDB struct {
	Users     []string
	MailBoxes map[string]*ImapMailBox
}

func (m MemDB) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewMemDB() *MemDB {
	conf := cfg.CurBackendConf()
	memDbFilePath := path.Join(conf.CurDBHome, string(filepath.Separator), util.DBMemHome, string(filepath.Separator), MemDbFileName)
	md := &MemDB{}
	err := util.ReadJsonFile(memDbFilePath, md)
	if err != nil {
		panic(err)
	}
	return md
}

func InitDefaultData(memHome string) (err error) {
	err = util.TouchDir(memHome)
	if err != nil {
		return err
	}

	mDBPath := path.Join(memHome, string(filepath.Separator), MemDbFileName)

	mdb := &MemDB{}
	mdb.Users = defaultUser
	mdb.MailBoxes = defaultMailBox
	err = util.WriteJsonFile(mDBPath, mdb)
	if err != nil {
		return err
	}
	return nil
}
