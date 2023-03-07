package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"sync"
)

const (
	MemDbFileName = "memDB.json"
)

var (
	_memBackLog = logrus.WithFields(logrus.Fields{
		"mode": "backend storage",
	})
	defaultUser = map[string]string{
		"ribencong@simplenets.org": "123",
		"fireman@simplenets.org":   "123",
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
	sync.RWMutex
	imapUsers   map[string]*ImapUser
	Users       map[string]string
	MailBoxes   MailBoxGroup
	msgId       uint32
	smtpSession map[string]*SmtpSession
}

func NewMemDB() *MemDB {
	conf := cfg.CurBackendConf()
	memDbFilePath := path.Join(conf.CurDBHome, string(filepath.Separator),
		util.DBMemHome, string(filepath.Separator), MemDbFileName)
	md := &MemDB{
		imapUsers:   make(map[string]*ImapUser),
		Users:       make(map[string]string),
		MailBoxes:   make(MailBoxGroup),
		smtpSession: make(map[string]*SmtpSession),
	}
	err := util.ReadJsonFile(memDbFilePath, md)
	if err != nil {
		panic(err)
	}
	for username, password := range md.Users {
		var user = NewUser(username, password, md)
		md.imapUsers[username] = user
	}
	return md
}

func NewUser(user, pwd string, de *MemDB) *ImapUser {
	iu := &ImapUser{
		username: user,
		password: pwd,
		db:       de,
	}
	iu.mailboxes = de.MailBoxGrp(iu)
	return iu
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
