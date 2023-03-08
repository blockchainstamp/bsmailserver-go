package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/backstore"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/emersion/go-smtp"
	"time"
)

type SmtpSrv struct {
	service *smtp.Server
}

func newSmtpSrv() *SmtpSrv {

	var ss = &SmtpSrv{}
	var conf = cfg.CurSmtpConf()
	var smtpSrv = smtp.NewServer(backstore.Inst())

	smtpSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	smtpSrv.Domain = conf.SrvDomain
	smtpSrv.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	smtpSrv.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	smtpSrv.MaxMessageBytes = conf.MaxMessageBytes
	smtpSrv.MaxRecipients = conf.MaxRecipients
	smtpSrv.AllowInsecureAuth = conf.TlsCfg == nil
	smtpSrv.TLSConfig = conf.TlsCfg

	ss.service = smtpSrv
	_srvLog.Info("smtp service init success at:", smtpSrv.Addr)

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
				_srvLog.Error("imap service exit:", err)
			}
		} else {
			err := ss.service.ListenAndServeTLS()
			if err != nil {
				_srvLog.Error("imap service exit:", err)
			}
		}
	}()
}
