package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

func (m *MemDB) LogoutFromImap(user *ImapUser) {

}

func (m *MemDB) MailBoxGrp(user *ImapUser) MailBoxGroup {
	group := make(MailBoxGroup)
	for s, box := range m.MailBoxes {
		group[s] = box
		group[s].User = user
	}
	return group
}

func (m *MemDB) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	m.Lock()
	defer m.Unlock()
	pwd, ok := m.Users[username]
	if !ok {
		return nil, util.BackendNoSuchUser
	}
	if pwd != password {
		return nil, util.BackendBadUser
	}
	user, ok := m.imapUsers[username]
	if ok {
		return user, nil
	}
	user = NewUser(username, password, m)
	m.imapUsers[username] = user
	return user, nil
}
