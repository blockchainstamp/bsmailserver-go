package protocol

import (
	"github.com/emersion/go-imap"
	"time"
)

type ImapMailBox struct {
	Subscribed bool
	BName      string
	BInfo      *imap.MailboxInfo
}

func (mailbox *ImapMailBox) Name() string {
	_protLog.Debugf("[%s] Name ", mailbox.BName)
	return mailbox.BName
}

func (mailbox *ImapMailBox) Info() (*imap.MailboxInfo, error) {
	_protLog.Debugf("[%s] Info ", mailbox.BName)
	return mailbox.BInfo, nil
}

func (mailbox *ImapMailBox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	_protLog.Debugf("[%s] Status %v", mailbox.BName, items)

	return imap.NewMailboxStatus(mailbox.BName, items), nil
}

func (mailbox *ImapMailBox) SetSubscribed(subscribed bool) error {
	_protLog.Debugf("[%s] SetSubscribed ", mailbox.BName)
	mailbox.Subscribed = subscribed
	return nil
}

func (mailbox *ImapMailBox) Check() error {
	_protLog.Debugf("[%s]Mailbox Check", mailbox.BName)
	return nil
}

func (mailbox *ImapMailBox) ListMessages(uid bool, seqset *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	_protLog.Debugf("[%s]Mailbox ListMessages", mailbox.BName)
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
	_protLog.Debugf("[%s]Mailbox SearchMessages", mailbox.BName)
	if uid {
		return nil, nil
	} else {
		return nil, nil
	}
}

func (mailbox *ImapMailBox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	_protLog.Debugf("[%s]Mailbox CreateMessage", mailbox.BName)
	return nil
}

func (mailbox *ImapMailBox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	_protLog.Debugf("[%s]Mailbox UpdateMessagesFlags", mailbox.BName)
	return nil
}

func (mailbox *ImapMailBox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	_protLog.Debugf("[%s]Mailbox CopyMessages", mailbox.BName)
	return nil
}

func (mailbox *ImapMailBox) Expunge() error {
	_protLog.Debugf("[%s]Mailbox Expunge", mailbox.BName)
	return nil
}
