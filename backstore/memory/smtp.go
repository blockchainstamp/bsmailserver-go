package memory

import (
	"bytes"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-msgauth/dkim"
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
	data     *bytes.Buffer
}

func (ss *SmtpSession) Reset() {

}

func (ss *SmtpSession) Logout() error {
	ss.finalize(ss)
	return nil
}

func (ss *SmtpSession) AuthPlain(username, password string) error {
	_memBackLog.Info("auth plain for:", username)
	return ss.auth(username, password)
}

func (ss *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	_memBackLog.Debug("mail from:", from, opts)
	ss.from = from
	return nil
}

func (ss *SmtpSession) Rcpt(to string) error {
	_memBackLog.Debug("to:", to)
	ss.tos = append(ss.tos, to)
	return nil
}

func (ss *SmtpSession) Data(r io.Reader) error {
	if ss.from != "" {
		_, err := ss.data.ReadFrom(r)
		return err
	}
	var buf bytes.Buffer
	tr := io.TeeReader(r, &buf)
	_, err := ss.data.ReadFrom(tr)
	if err != nil {
		return err
	}
	vers, err := dkim.Verify(&buf)
	if err != nil {
		return err
	}
	for _, ver := range vers {
		if ver.Err != nil {
			_memBackLog.Warn("dkim sig err:", ver.Domain, ver.Err)
			return err
		} else {
			_memBackLog.Info("dkim sig success:", ver.Domain)
		}
	}

	return err
}

func (m *MemDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	atomic.AddUint32(&m.msgId, 1)
	return &SmtpSession{
		id:       m.msgId,
		auth:     m.authSession,
		finalize: m.finalizeSession,
		data:     new(bytes.Buffer),
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
			_memBackLog.Warn("outer user now:")
			continue
		}
		mbox, err := uname.GetMailbox("INBOX")
		if err != nil {
			_memBackLog.Warn("mail box err:", err)
			continue
		}
		err = mbox.CreateMessage([]string{}, time.Now(), s.data)
		if err != nil {
			_memBackLog.Warn("create msg err:", err)
			continue
		}
	}
}
