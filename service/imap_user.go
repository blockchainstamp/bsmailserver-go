package service

import (
	"github.com/emersion/go-imap/backend"
)

type ImapUser struct {
	username  string
	password  string
	mailboxes map[string]*ImapMailBox
}

func (user *ImapUser) Username() string {
	return user.username
}

func (user *ImapUser) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	_imapLog.Debugf("[%s] ListMailboxes[%t]", user.username, subscribed)
	var boxes []backend.Mailbox
	for _, box := range user.mailboxes {
		boxes = append(boxes, box)
	}
	return boxes, nil
}

func (user *ImapUser) GetMailbox(name string) (backend.Mailbox, error) {
	_imapLog.Debugf("[%s] GetMailbox[%s]", user.username, name)
	return user.mailboxes[name], nil
}

func (user *ImapUser) CreateMailbox(name string) error {
	_imapLog.Debugf("[%s] CreateMailbox", user.username)
	return nil
}

func (user *ImapUser) DeleteMailbox(name string) error {
	_imapLog.Debugf("[%s] DeleteMailbox", user.username)
	return nil
}

func (user *ImapUser) RenameMailbox(existingName, newName string) error {

	_imapLog.Debugf("[%s] RenameMailbox", user.username)
	return nil
}

func (user *ImapUser) Logout() error {
	_imapLog.Debugf("[%s] Logout", user.username)
	return nil
}
