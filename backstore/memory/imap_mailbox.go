package memory

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/backend/backendutil"
	"io"
	"time"
)

type ImapMailBox struct {
	Subscribed bool
	BName      string
	BInfo      *imap.MailboxInfo
	Messages   MessageList
	User       *ImapUser
}

func (mailbox *ImapMailBox) Name() string {
	_memBackLog.Debugf("[%s] Name ", mailbox.BName)
	return mailbox.BName
}

func (mailbox *ImapMailBox) Info() (*imap.MailboxInfo, error) {
	_memBackLog.Debugf("[%s] Info ", mailbox.BName)
	return mailbox.BInfo, nil
}

func (mailbox *ImapMailBox) uidNext() uint32 {
	var uid uint32
	for _, msg := range mailbox.Messages {
		if msg.Uid > uid {
			uid = msg.Uid
		}
	}
	uid++
	return uid
}

func (mailbox *ImapMailBox) flags() []string {
	flagsMap := make(map[string]bool)
	for _, msg := range mailbox.Messages {
		for _, f := range msg.Flags {
			if !flagsMap[f] {
				flagsMap[f] = true
			}
		}
	}

	var flags []string
	for f := range flagsMap {
		flags = append(flags, f)
	}
	return flags
}

func (mailbox *ImapMailBox) unseenSeqNum() uint32 {
	for i, msg := range mailbox.Messages {
		seqNum := uint32(i + 1)

		seen := false
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				seen = true
				break
			}
		}

		if !seen {
			return seqNum
		}
	}
	return 0
}

func (mailbox *ImapMailBox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	_memBackLog.Debugf("[%s] Status %v", mailbox.BName, items)

	status := imap.NewMailboxStatus(mailbox.BName, items)
	status.Flags = mailbox.flags()
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = mailbox.unseenSeqNum()

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = uint32(len(mailbox.Messages))
		case imap.StatusUidNext:
			status.UidNext = mailbox.uidNext()
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = 0 // TODO
		case imap.StatusUnseen:
			status.Unseen = 0 // TODO
		}
	}

	return status, nil
}

func (mailbox *ImapMailBox) SetSubscribed(subscribed bool) error {
	_memBackLog.Debugf("[%s] SetSubscribed ", mailbox.BName)
	mailbox.Subscribed = subscribed
	return nil
}

func (mailbox *ImapMailBox) Check() error {
	_memBackLog.Debugf("[%s]Mailbox Check", mailbox.BName)
	return nil
}

func (mailbox *ImapMailBox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	_memBackLog.Debugf("[%s]Mailbox ListMessages", mailbox.BName)
	defer close(ch)

	for i, msg := range mailbox.Messages {
		seqNum := uint32(i + 1)

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		if !seqSet.Contains(id) {
			continue
		}

		m, err := msg.Fetch(seqNum, items)
		if err != nil {
			continue
		}
		ch <- m
	}

	return nil
}

func (mailbox *ImapMailBox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	_memBackLog.Debugf("[%s]Mailbox SearchMessages", mailbox.BName)
	var ids []uint32
	for i, msg := range mailbox.Messages {
		seqNum := uint32(i + 1)

		ok, err := msg.Match(seqNum, criteria)
		if err != nil || !ok {
			continue
		}

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (mailbox *ImapMailBox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	_memBackLog.Debugf("[%s]Mailbox CreateMessage", mailbox.BName)
	if date.IsZero() {
		date = time.Now()
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	mailbox.Messages = append(mailbox.Messages, &Message{
		Uid:   mailbox.uidNext(),
		Date:  date,
		Size:  uint32(len(b)),
		Flags: flags,
		Body:  b,
	})
	return nil
}

func (mailbox *ImapMailBox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	_memBackLog.Debugf("[%s]Mailbox UpdateMessagesFlags", mailbox.BName)
	for i, msg := range mailbox.Messages {
		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		msg.Flags = backendutil.UpdateFlags(msg.Flags, operation, flags)
	}

	return nil
}

func (mailbox *ImapMailBox) CopyMessages(uid bool, seqset *imap.SeqSet, destName string) error {
	_memBackLog.Debugf("[%s]Mailbox CopyMessages", mailbox.BName)
	dest, ok := mailbox.User.mailboxes[destName]
	if !ok {
		return backend.ErrNoSuchMailbox
	}

	for i, msg := range mailbox.Messages {
		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		msgCopy := *msg
		msgCopy.Uid = dest.uidNext()
		dest.Messages = append(dest.Messages, &msgCopy)
	}

	return nil
}

func (mailbox *ImapMailBox) Expunge() error {
	_memBackLog.Debugf("[%s]Mailbox Expunge", mailbox.BName)
	for i := len(mailbox.Messages) - 1; i >= 0; i-- {
		msg := mailbox.Messages[i]

		deleted := false
		for _, flag := range msg.Flags {
			if flag == imap.DeletedFlag {
				deleted = true
				break
			}
		}

		if deleted {
			mailbox.Messages = append(mailbox.Messages[:i], mailbox.Messages[i+1:]...)
		}
	}

	return nil
}
