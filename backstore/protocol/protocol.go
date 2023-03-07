package protocol

import "github.com/sirupsen/logrus"

var (
	_protLog = logrus.WithFields(logrus.Fields{
		"mode": "backend storage",
	})
)

type DatabaseExe interface {
	MailBoxGrp(user *ImapUser) MailBoxGroup
	Logout(user *ImapUser)
}
type MailBoxGroup map[string]*ImapMailBox

func (mbg MailBoxGroup) Copy() MailBoxGroup {
	newGrp := make(MailBoxGroup)
	for s, box := range mbg {
		newGrp[s] = box
	}
	return newGrp
}

func NewUser(user, pwd string, de DatabaseExe) *ImapUser {
	iu := &ImapUser{
		username: user,
		password: pwd,
		dbExe:    de,
	}
	iu.mailboxes = de.MailBoxGrp(iu)
	return iu
}
