package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/backstore/protocol"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
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
	defaultMailBox = map[string]*protocol.ImapMailBox{
		"INBOX": &protocol.ImapMailBox{
			Subscribed: false,
			BName:      "INBOX",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{},
				Delimiter:  "/",
				Name:       "INBOX",
			},
		},
		"草稿箱": &protocol.ImapMailBox{
			Subscribed: false,
			BName:      "草稿箱",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Drafts"},
				Delimiter:  "/",
				Name:       "草稿箱",
			},
		},
		"已发送": &protocol.ImapMailBox{
			Subscribed: false,
			BName:      "已发送",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Sent"},
				Delimiter:  "/",
				Name:       "已发送",
			},
		},
		"已删除": &protocol.ImapMailBox{
			Subscribed: false,
			BName:      "已删除",
			BInfo: &imap.MailboxInfo{
				Attributes: []string{"\\Trash"},
				Delimiter:  "/",
				Name:       "已删除",
			},
		},
		"区块链邮票": &protocol.ImapMailBox{
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
	userOnline map[string]*protocol.ImapUser
	Users      map[string]string
	MailBoxes  protocol.MailBoxGroup
}

func (m *MemDB) Logout(user *protocol.ImapUser) {
	m.Lock()
	defer m.Unlock()
	delete(m.Users, user.Username())
}

func (m *MemDB) MailBoxGrp(*protocol.ImapUser) protocol.MailBoxGroup {
	return m.MailBoxes
}

func (m *MemDB) MailBoxesOf(user *protocol.ImapUser) protocol.MailBoxGroup {
	group := make(protocol.MailBoxGroup)
	for s, box := range m.MailBoxes {
		group[s] = box
	}
	return group
}

func (m *MemDB) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	m.Lock()
	defer m.Unlock()
	pwd, ok := m.Users[username]
	if !ok {
		return nil, util.ImapNoSuchUser
	}
	if pwd != password {
		return nil, util.ImapBadUser
	}
	user, ok := m.userOnline[username]
	if ok {
		return user, nil
	}
	user = protocol.NewUser(username, password, m)
	m.userOnline[username] = user
	return user, nil
}

func NewMemDB() *MemDB {
	conf := cfg.CurBackendConf()
	memDbFilePath := path.Join(conf.CurDBHome, string(filepath.Separator),
		util.DBMemHome, string(filepath.Separator), MemDbFileName)
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
