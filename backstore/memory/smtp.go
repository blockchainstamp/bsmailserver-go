package memory

import (
	"bytes"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	"io"
	"sync/atomic"
	"time"
)

type SmtpEnvelop struct {
	From       string
	Tos        []string
	EmailID    string
	DataSource io.Reader
}
type SmtpSession struct {
	id          uint32
	isLcl       bool
	auth        func(string, string) error
	serveForCli func(env *SmtpEnvelop) error
	serveForSrv func(env *SmtpEnvelop) error
	ev          *SmtpEnvelop
	db          *MemDB
}

func (ss *SmtpSession) Reset() {
}

func (ss *SmtpSession) Logout() error {
	return nil
}

func (ss *SmtpSession) AuthPlain(username, password string) error {
	_memBackLog.Info("auth plain for:", username)
	err := ss.auth(username, password)
	if err != nil {
		_memBackLog.Warn("no such user: ", username)
		return err
	}
	ss.isLcl = true
	return nil
}

func (ss *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	_memBackLog.Debug("mail from:", from, opts)
	ss.ev.From = from
	return nil
}

func (ss *SmtpSession) Rcpt(to string) error {
	_memBackLog.Debug("to: ", to)
	if ss.isLcl {
		ss.ev.Tos = append(ss.ev.Tos, to)
		return nil
	}
	_, ok := ss.db.Users[to]
	if !ok {
		_memBackLog.Warn("no such user ", to)
		return util.BackendNoSuchUser
	}
	ss.ev.Tos = append(ss.ev.Tos, to)
	return nil
}

func (ss *SmtpSession) Data(r io.Reader) error {
	_memBackLog.Debugf("local[%t], process mail data for[%s]", ss.isLcl, ss.ev.From)
	ss.ev.DataSource = r
	//TODO:: get stamp here
	if ss.isLcl {
		return ss.serveForCli(ss.ev)
	}
	return ss.serveForSrv(ss.ev)

}

func (m *MemDB) NewSession(c *smtp.Conn) (smtp.Session, error) {
	atomic.AddUint32(&m.sessionID, 1)
	_memBackLog.Debug("new session create: ", m.sessionID)
	return &SmtpSession{
		id:          m.sessionID,
		ev:          &SmtpEnvelop{},
		auth:        m.authSession,
		serveForSrv: m.serveForSrv,
		serveForCli: m.serveForCli,
		db:          m,
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

func (m *MemDB) serveForCli(env *SmtpEnvelop) error {
	var outerTos = make(map[string][]string)
	for _, to := range env.Tos {
		user, ok := m.imapUsers[to]
		if !ok {
			_, domain, err := util.ParseEmailAddress(to)
			if err != nil {
				_memBackLog.Warn("invalid to from client: ", to, err)
				return err
			}
			outerTos[domain] = append(outerTos[domain], to)
			continue
		}
		if err := m.msgToInbox(user, env); err != nil {
			_memBackLog.Warn("create client msg err: ", to, err)
			return err
		}
	}
	if len(outerTos) == 0 {
		_memBackLog.Debug("no outer tos:", env.Tos)
		return nil
	}

	return m.sendToRemoteSrv(outerTos, env)
}

func (m *MemDB) serveForSrv(env *SmtpEnvelop) error {
	var buf bytes.Buffer
	tr := io.TeeReader(env.DataSource, &buf)
	verifications, err := dkim.Verify(tr)
	if err != nil {
		return err
	}
	for _, ver := range verifications {
		if ver.Err != nil {
			_memBackLog.Warn("dkim sig err:", ver.Domain, ver.Err)
			return err
		} else {
			_memBackLog.Info("dkim sig success:", ver.Domain)
		}
	}
	for _, to := range env.Tos {
		user, ok := m.imapUsers[to]
		if !ok {
			return util.BackendNoSuchUser
		}
		if err := m.msgToInbox(user, env); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemDB) msgToInbox(user *ImapUser, env *SmtpEnvelop) error {
	inbox, err := user.GetMailbox(util.MBXIndex)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(env.DataSource)
	if err != nil {
		return err
	}
	return inbox.CreateMessage([]string{}, time.Now(), &buf)
}

func sendMail(cli *smtp.Client, tos []string, from string, envData io.Reader) error {
	err := cli.Mail(from, nil)
	if err != nil {
		_memBackLog.Warn("mail err:", err)
		return err
	}
	for _, to := range tos {
		err = cli.Rcpt(to)
		if err != nil {
			_memBackLog.Warn("rcpt err:", to, err)
			return err
		}
	}
	wc, err := cli.Data()
	if err != nil {
		_memBackLog.Warn("data err:", err)
		return err
	}

	_, err = io.Copy(wc, envData)
	if err != nil {
		_memBackLog.Warn("write to err:", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		_memBackLog.Warn("close err:", err)
		return err
	}
	err = cli.Quit()
	if err != nil {
		_memBackLog.Warn("quit err:", err)
		return err
	}
	_memBackLog.Info("send mail success: ", from, tos)
	return nil
}

func (m *MemDB) sendToRemoteSrv(tos map[string][]string, env *SmtpEnvelop) error {
	conf := cfg.CurSmtpConf()

	for domain, subTos := range tos {
		_memBackLog.Debug("prepare to send out mail: ", domain, subTos)
		cli, err := util.DnsInst().ValidSmtpCli(conf.SrvDomain, domain)
		if err != nil {
			_memBackLog.Warn("no valid mail cli: ", err)
			return err
		}
		options := &dkim.SignOptions{
			Domain:   conf.SrvDomain,
			Selector: conf.DKIMSelector,
			Signer:   conf.DKIMSigner,
		}

		var buff bytes.Buffer
		if err := dkim.Sign(&buff, env.DataSource, options); err != nil {
			_memBackLog.Warn("sign for mail failed: ", err, env.From)
			return err
		}
		err = sendMail(cli, subTos, env.From, &buff)
		if err != nil {
			_memBackLog.Warn("send out mail failed: ", err, env.From, subTos)
			return err
		}
		_ = cli.Close()
	}
	return nil
}
