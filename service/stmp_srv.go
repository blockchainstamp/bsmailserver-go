package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/emersion/go-smtp"
	"io"
	"time"
)

func (ss *SmtpSrv) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return ss, nil
}

type SmtpSrv struct {
	service *smtp.Server
}

func (ss *SmtpSrv) Reset() {
	_smtpLog.Debug("smtp reset")
}

func (ss *SmtpSrv) Logout() error {
	_smtpLog.Debug("smtp Logout")
	return nil
}

func (ss *SmtpSrv) AuthPlain(username, password string) error {
	_smtpLog.Debug("smtp AuthPlain", username)
	return nil
}

func (ss *SmtpSrv) Mail(from string, opts *smtp.MailOptions) error {
	_smtpLog.Debug("smtp Mail", from)
	return nil
}

func (ss *SmtpSrv) Rcpt(to string) error {
	_smtpLog.Debug("smtp Rcpt", to)
	return nil
}

func (ss *SmtpSrv) Data(r io.Reader) error {
	_smtpLog.Debug("smtp Data")
	return nil
}

func newSmtpSrv() *SmtpSrv {

	var ss = &SmtpSrv{}
	var conf = cfg.CurSmtpConf()
	var smtpSrv = smtp.NewServer(ss)

	smtpSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	smtpSrv.Domain = conf.SrvDomain
	smtpSrv.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	smtpSrv.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	smtpSrv.MaxMessageBytes = conf.MaxMessageBytes
	smtpSrv.MaxRecipients = conf.MaxRecipients
	smtpSrv.AllowInsecureAuth = conf.TlsCfg == nil
	smtpSrv.TLSConfig = conf.TlsCfg

	ss.service = smtpSrv
	_smtpLog.Info("smtp service init success at:", smtpSrv.Addr)

	return ss
}

func (ss *SmtpSrv) Close() error {
	return ss.service.Close()
}

func (ss *SmtpSrv) Run() {
	go func() {
		var conf = cfg.CurSmtpConf()
		if conf.AllowNotSecure {
			err := ss.service.ListenAndServe()
			if err != nil {
				panic(err)
			}
		} else {
			err := ss.service.ListenAndServeTLS()
			if err != nil {
				panic(err)
			}
		}
	}()
}
