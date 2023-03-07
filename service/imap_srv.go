package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/backstore"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	iSrv "github.com/emersion/go-imap/server"
)

type ImapSrv struct {
	imapSrv *iSrv.Server
}

func newImapSrv() *ImapSrv {
	var is = &ImapSrv{}
	var conf = cfg.CurImapConf()
	imapSrv := iSrv.New(backstore.Inst())
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
