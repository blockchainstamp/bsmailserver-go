package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	iSrv "github.com/emersion/go-imap/server"
)

type ImapSrv struct {
	imapSrv *iSrv.Server
}

func (ic *ImapSrv) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	u := &ImapUser{username: username, password: password}
	_imapLog.Infof("user[%s] login success", username)
	return u, nil
}

func newImapSrv() *ImapSrv {
	var ic = &ImapSrv{}
	var conf = cfg.CurImapConf()
	imapSrv := iSrv.New(ic)
	imapSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	imapSrv.AllowInsecureAuth = conf.TlsCfg == nil
	imapSrv.TLSConfig = conf.TlsCfg
	ic.imapSrv = imapSrv
	_smtpLog.Info("smtp service init success at:", imapSrv.Addr)
	return ic
}

func (ic *ImapSrv) Close() error {
	return ic.imapSrv.Close()
}

func (ic *ImapSrv) Run() {
	go func() {
		var conf = cfg.CurImapConf()
		if conf.AllowNotSecure {
			err := ic.imapSrv.ListenAndServe()
			if err != nil {
				panic(err)
			}
		} else {
			err := ic.imapSrv.ListenAndServeTLS()
			if err != nil {
				panic(err)
			}
		}
	}()
}
