package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-imap/backend"
)

type ImapUser struct {
	username  string
	password  string
	mailboxes MailBoxGroup
	db        *MemDB
}

func (user *ImapUser) Username() string {
	return user.username
}

func (user *ImapUser) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	_memBackLog.Debugf("[%s] ListMailboxes[%t]", user.username, subscribed)
	var boxes []backend.Mailbox
	for _, box := range user.mailboxes {
		boxes = append(boxes, box)
	}
	return boxes, nil
}

func (user *ImapUser) GetMailbox(name string) (backend.Mailbox, error) {
	_memBackLog.Debugf("[%s] GetMailbox[%s]", user.username, name)
	mbox, ok := user.mailboxes[name]
	if !ok {
		return nil, util.BackendNoSuchMailBox
	}

	return mbox, nil
}

func (user *ImapUser) CreateMailbox(name string) error {
	_memBackLog.Debugf("[%s] CreateMailbox", user.username)
	if _, ok := user.mailboxes[name]; ok {
		return util.BackendDuplicateMailBox
	}
	user.mailboxes[name] = &ImapMailBox{
		BName: name,
	}
	//TODO:: user.dbExe.CreateMailbox
	return nil
}

func (user *ImapUser) DeleteMailbox(name string) error {
	_memBackLog.Debugf("[%s] DeleteMailbox", user.username)
	if name == "INBOX" {
		return util.BackendDeleteInbox
	}
	if _, ok := user.mailboxes[name]; !ok {
		return util.BackendNoSuchMailBox
	}
	delete(user.mailboxes, name)
	//TODO:: user.dbExe.DeleteMailbox
	return nil
}

func (user *ImapUser) RenameMailbox(existingName, newName string) error {
	_memBackLog.Debugf("[%s] RenameMailbox", user.username)
	mbox, ok := user.mailboxes[existingName]
	if !ok {
		return util.BackendNoSuchMailBox
	}

	user.mailboxes[newName] = &ImapMailBox{
		BName:    newName,
		Messages: mbox.Messages,
	}

	mbox.Messages = nil

	if existingName != "INBOX" {
		delete(user.mailboxes, existingName)
	}
	//TODO:: user.dbExe.RenameMailbox
	return nil
}

func (user *ImapUser) Logout() error {
	_memBackLog.Debugf("[%s] Logout", user.username)
	user.db.LogoutFromImap(user)
	return nil
}

type MailBoxGroup map[string]*ImapMailBox

func (mbg MailBoxGroup) Copy() MailBoxGroup {
	newGrp := make(MailBoxGroup)
	for s, box := range mbg {
		newGrp[s] = box
	}
	return newGrp
}

type MessageList []*Message
