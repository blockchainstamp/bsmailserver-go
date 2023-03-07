package memory

import (
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-smtp"
	"io"
	"sync/atomic"
	"time"
)

type SmtpSession struct {
	id       uint32
	from     string
	tos      []string
	auth     func(string, string) error
	finalize func(session *SmtpSession)
	reader   io.Reader
}

func (ss *SmtpSession) Reset() {

}

func (ss *SmtpSession) Logout() error {
	ss.finalize(ss)
	return nil
}

func (ss *SmtpSession) AuthPlain(username, password string) error {
	return ss.auth(username, password)
}

func (ss *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	ss.from = from
	return nil
}

func (ss *SmtpSession) Rcpt(to string) error {
	ss.tos = append(ss.tos, to)
	return nil
}

func (ss *SmtpSession) Data(r io.Reader) error {
	ss.reader = r
	return nil
}

func (ss *SmtpSession) Len() int {
	return 0
}
func (ss *SmtpSession) Read(p []byte) (n int, err error) {
	return ss.reader.Read(p)
}
func (m *MemDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	atomic.AddUint32(&m.msgId, 1)
	return &SmtpSession{
		id:       m.msgId,
		auth:     m.authSession,
		finalize: m.finalizeSession,
	}, nil
}

func (m *MemDB) authSession(username, password string) error {
	pwd, ok := m.Users[username]
	if !ok {
		return util.BackendNoSuchUser
	}
	if pwd != password {
		return util.BackendBadUser
	}
	return nil
}

func (m *MemDB) finalizeSession(s *SmtpSession) {

	for _, to := range s.tos {
		uname, ok := m.imapUsers[to]
		if !ok {
			_memBackLog.Warn("unSupport user now:", to)
			continue
		}
		mbox, err := uname.GetMailbox("INBOX")
		if err != nil {
			_memBackLog.Warn("mail box err:", err)
			continue
		}
		err = mbox.CreateMessage([]string{"\\Seen"}, time.Now(), s)
		if err != nil {
			_memBackLog.Warn("create msg err:", err)
			continue
		}
	}
}
