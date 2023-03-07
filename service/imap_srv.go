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

func (is *ImapSrv) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {

	u := &ImapUser{username: username, password: password}
	_imapLog.Infof("user[%s] login success", username)
	return u, nil
}

func newImapSrv() *ImapSrv {
	var is = &ImapSrv{}
	var conf = cfg.CurImapConf()
	imapSrv := iSrv.New(is)
	imapSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	imapSrv.AllowInsecureAuth = conf.TlsCfg == nil
	imapSrv.TLSConfig = conf.TlsCfg
	is.imapSrv = imapSrv
	_smtpLog.Info("smtp service init success at:", imapSrv.Addr)
	return is
}

func (is *ImapSrv) Close() error {
	return is.imapSrv.Close()
}

func (is *ImapSrv) Run() {
	go func() {
		var conf = cfg.CurImapConf()
		if conf.AllowNotSecure {
			err := is.imapSrv.ListenAndServe()
			if err != nil {
				panic(err)
			}
		} else {
			err := is.imapSrv.ListenAndServeTLS()
			if err != nil {
				panic(err)
			}
		}
	}()
}
