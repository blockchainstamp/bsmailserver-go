package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/backstore"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-smtp"
	"time"
)

type SmtpSrv struct {
	service *smtp.Server
	smtp25  *smtp.Server
}

func newSmtpSrv() *SmtpSrv {

	var ss = &SmtpSrv{}
	var conf = cfg.CurSmtpConf()
	var smtpSrv = smtp.NewServer(backstore.Inst())
	var smtp25Srv = smtp.NewServer(backstore.Inst())

	smtpSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	smtpSrv.Domain = conf.SrvDomain
	smtpSrv.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	smtpSrv.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	smtpSrv.MaxMessageBytes = conf.MaxMessageBytes
	smtpSrv.MaxRecipients = conf.MaxRecipients
	smtpSrv.AllowInsecureAuth = conf.AllowNotSecure
	smtpSrv.TLSConfig = conf.TlsCfg

	ss.service = smtpSrv
	_srvLog.Info("smtp service init success at: ", smtpSrv.Addr)

	smtp25Srv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, util.DefaultSystemSmtpPort)
	smtp25Srv.Domain = conf.SrvDomain
	smtp25Srv.ReadTimeout = time.Duration(conf.ReadTimeOut) * time.Second
	smtp25Srv.WriteTimeout = time.Duration(conf.WriteTimeOut) * time.Second
	smtp25Srv.MaxMessageBytes = conf.MaxMessageBytes
	smtp25Srv.MaxRecipients = conf.MaxRecipients
	smtp25Srv.AllowInsecureAuth = conf.AllowNotSecure
	smtp25Srv.TLSConfig = conf.TlsCfg

	ss.smtp25 = smtp25Srv
	_srvLog.Info("default smtp service init success at: ", smtp25Srv.Addr)

	return ss
}

func (ss *SmtpSrv) Close() error {
	err1, err2 := ss.service.Close(), ss.smtp25.Close()
	if err2 == nil && err1 == nil {
		return nil
	}
	return fmt.Errorf("%v:%v", err1, err2)
}

func (ss *SmtpSrv) Run() {
	go func() {
		_srvLog.Info("starting smtp service.......")
		err := ss.service.ListenAndServeTLS()
		if err != nil {
			_srvLog.Error("smtp service exit:", err)
			panic(err)
		}
	}()
	go func() {
		_srvLog.Info("starting system 25 smtp service.......")
		err := ss.smtp25.ListenAndServeTLS()
		if err != nil {
			_srvLog.Error("system 25 smtp service exit:", err)
			panic(err)
		}
	}()
}
