package service

import (
	"github.com/emersion/go-imap"
	"time"
)

type ImapMailBox struct {
	Subscribed bool
	name       string
	user       *ImapUser
	info       *imap.MailboxInfo
}

func (mailbox *ImapMailBox) Name() string {
	_imapLog.Debugf("[%s] Name ", mailbox.name)
	return mailbox.name
}

func (mailbox *ImapMailBox) Info() (*imap.MailboxInfo, error) {
	_imapLog.Debugf("[%s] Info ", mailbox.name)
	return mailbox.info, nil
}

func (mailbox *ImapMailBox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	_imapLog.Debugf("[%s] Status %v", mailbox.name, items)

	return imap.NewMailboxStatus(mailbox.name, items), nil
}

func (mailbox *ImapMailBox) SetSubscribed(subscribed bool) error {
	_imapLog.Debugf("[%s] SetSubscribed ", mailbox.name)
	mailbox.Subscribed = subscribed
	return nil
}

func (mailbox *ImapMailBox) Check() error {
	_imapLog.Debugf("[%s]Mailbox Check", mailbox.name)
	return nil
}

func (mailbox *ImapMailBox) ListMessages(uid bool, seqset *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	_imapLog.Debugf("[%s]Mailbox ListMessages", mailbox.name)
	messages := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		if uid {
			//done <- mailbox.user.cli.UidFetch(seqSet, items, messages)
		} else {
			//done <- mailbox.user.cli.Fetch(seqSet, items, messages)
		}
		done <- nil
	}()

	for msg := range messages {
		ch <- msg
	}
	return <-done
}

func (mailbox *ImapMailBox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	_imapLog.Debugf("[%s]Mailbox SearchMessages", mailbox.name)
	if uid {
		return nil, nil
	} else {
		return nil, nil
	}
}

func (mailbox *ImapMailBox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	_imapLog.Debugf("[%s]Mailbox CreateMessage", mailbox.name)
	return nil
}

func (mailbox *ImapMailBox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	_imapLog.Debugf("[%s]Mailbox UpdateMessagesFlags", mailbox.name)
	return nil
}

func (mailbox *ImapMailBox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	_imapLog.Debugf("[%s]Mailbox CopyMessages", mailbox.name)
	return nil
}

func (mailbox *ImapMailBox) Expunge() error {
	_imapLog.Debugf("[%s]Mailbox Expunge", mailbox.name)
	return nil
}
